聚合网关实践
let vauleMaps = Object.fromEntries(metadata['helloworld.GetUserInfoReply']);
retrunMsg['age'] = vauleMaps["helloworld.GetUserInfoReply"]["age"];retrunMsg['balance']= vauleMaps["helloworld.GetWalletInfoReply"]["balance"];