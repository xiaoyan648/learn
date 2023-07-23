package sort

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
)

// 1. 大文件排序 大文件-> 小临时文件 -> 内存归并排序
// 2.posting 很多小文件 -> 分组 -> 组织posting list -> posting list key排序存储 -> 每组一个文件 -> 多路归并排序
const (
	chunkSize = 200 // 每个块的大小
)

func MergeSort(filePath string) {
	// 打开原始文件和目标文件
	// filePath := "input.txt"
	inputFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// 分割文件为多个块并进行排序
	chunks, err := splitAndSortChunks(inputFile)
	if err != nil {
		panic(err)
	}

	// 合并块并写入最终的排序文件
	err = mergeChunks(chunks, outputFile)
	if err != nil {
		panic(err)
	}

	log.Println("排序完成！")
}

// 分割文件为多个块并进行排序
func splitAndSortChunks(inputFile *os.File) ([]string, error) {
	chunks := make([]string, 0)
	buffer := make([]int, 0, chunkSize)

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		i, _ := strconv.Atoi(scanner.Text())
		buffer = append(buffer, i)

		// 当缓冲区大小达到块大小时，进行排序并将块写入临时文件
		if len(buffer) >= chunkSize {
			sort.Ints(buffer)

			tmpFile, err := createTempFile()
			if err != nil {
				return nil, err
			}

			err = writeLinesToFile(buffer, tmpFile)
			if err != nil {
				return nil, err
			}

			chunks = append(chunks, tmpFile.Name())
			buffer = buffer[:0]
		}
	}

	// 处理剩余的数据
	if len(buffer) > 0 {
		sort.Ints(buffer)

		tmpFile, err := createTempFile()
		if err != nil {
			return nil, err
		}

		err = writeLinesToFile(buffer, tmpFile)
		if err != nil {
			return nil, err
		}

		chunks = append(chunks, tmpFile.Name())
	}

	return chunks, nil
}

// 创建临时文件
func createTempFile() (*os.File, error) {
	return os.CreateTemp("", "temp_chunk_*.txt")
}

// 将字符串切片逐行写入文件
func writeLinesToFile(lines []int, file *os.File) error {
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, line := range lines {
		_, err := writer.WriteString(strconv.Itoa(line) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// 合并块并写入最终的排序文件
// 这里没必要堆排序，直接for循环比较就可以了最后进行一次排序即可；堆排序 NlogN*K for循环 N*K + NlogN (M是小文件的数量)
// 堆排序适合直接大文件动态排序的场景，或者只需要topN的场景
func mergeChunks(chunks []string, outputFile *os.File) error {
	files := make([]*os.File, len(chunks))
	scanners := make([]*bufio.Scanner, len(chunks))

	// 打开临时文件并创建对应的扫描器
	for i, chunk := range chunks {
		file, err := os.Open(chunk)
		if err != nil {
			return err
		}
		defer file.Close()
		files[i] = file

		scanners[i] = bufio.NewScanner(file)
	}

	// 使用堆来进行多路归并排序
	var heapNodes HeapNodeList

	// 初始化堆
	for i, scanner := range scanners {
		if scanner.Scan() {
			value, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return err
			}
			heapNodes = append(heapNodes, &HeapNode{Value: value, ScannerIndex: i})
		} else {
			// 如果扫描器已经读取完毕，将其关闭并从堆中移除
			files[i].Close()
			log.Printf("文件 %s 已经读取完毕", chunks[i])
		}
	}
	// 建堆
	buildHeap(heapNodes)

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	// 多路归并排序
	for len(heapNodes) > 0 {
		// 从堆顶取出最小元素
		minNode := heapNodes[0]
		writer.WriteString(strconv.Itoa(minNode.Value) + "\n")

		// 从相应的扫描器中读取下一个元素
		scanner := scanners[minNode.ScannerIndex]
		if scanner.Scan() {
			value, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return err
			}
			minNode.Value = value
		} else {
			// 如果扫描器已经读取完毕，将其关闭并从堆中移除
			files[minNode.ScannerIndex].Close()
			heapNodes = removeHeapNode(heapNodes, 0)
		}

		// 调整堆以维持堆的性质
		heapify(heapNodes, 0)
	}

	return nil
}

// 堆结点
type HeapNode struct {
	Value        int // 结点值
	ScannerIndex int // 所属扫描器索引
}

// 堆结点列表
type HeapNodeList []*HeapNode

// 获取父结点索引
func getParentIndex(index int) int {
	return (index - 1) / 2
}

// 获取左子结点索引
func getLeftChildIndex(index int) int {
	return index*2 + 1
}

// 获取右子结点索引
func getRightChildIndex(index int) int {
	return index*2 + 2
}

// 建堆
func buildHeap(heap HeapNodeList) {
	lastNonLeafIndex := len(heap)/2 - 1
	for i := lastNonLeafIndex; i >= 0; i-- {
		heapify(heap, i)
	}
}

// getsmall

// 堆化
func heapify(heap HeapNodeList, index int) {
	smallest := index
	left := getLeftChildIndex(index)
	right := getRightChildIndex(index)

	if left < len(heap) && heap[left].Value < heap[smallest].Value {
		smallest = left
	}

	if right < len(heap) && heap[right].Value < heap[smallest].Value {
		smallest = right
	}

	if smallest != index {
		heap[index], heap[smallest] = heap[smallest], heap[index]
		heapify(heap, smallest)
	}
}

// 从堆中移除结点
func removeHeapNode(heap HeapNodeList, index int) HeapNodeList {
	heap[index] = heap[len(heap)-1]
	heap = heap[:len(heap)-1]
	heapify(heap, index)
	return heap
}
