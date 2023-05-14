package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRPC_eth_blockNumber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check the request parameters
		if r.Method != http.MethodPost {
			t.Errorf("expected %s request, got %s", http.MethodPost, r.Method)
		}
		if r.URL.Path != "/" {
			t.Errorf("expected request to path /, got %s", r.URL.Path)
		}

		// write the response
		response := `{"jsonrpc":"2.0","id":2,"result":"0x28bb63f"}`
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(response))
		if err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	defer server.Close()

	client := server.Client()
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"id":      2,
	}
	resp, err := makeRPCRequest(client, server.URL, reqBody)
	if err != nil {
		t.Fatalf("MakeRPCRequest returned unexpected error: %v", err)
	}

	expectedResp := `{"jsonrpc":"2.0","id":2,"result":"0x28bb63f"}`
	if string(resp) != expectedResp {
		t.Errorf("expected response %q, got %q", expectedResp, resp)
	}
}

// Meant to express how a mispell in method name would result if calling the Polygon RPC
func TestRPC_eth_blockNumberMistake(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check the request parameters
		if r.Method != http.MethodPost {
			t.Errorf("expected %s request, got %s", http.MethodPost, r.Method)
		}
		if r.URL.Path != "/" {
			t.Errorf("expected request to path /, got %s", r.URL.Path)
		}

		// parse the request body
		var requestBody struct {
			JSONRPC string      `json:"jsonrpc"`
			Method  string      `json:"method"`
			ID      interface{} `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("error decoding request body: %v", err)
		}

		// check the request method
		if requestBody.Method != "eth_blockNumberMistake" {
			// write an error response
			errorResponse := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      requestBody.ID,
				"error": map[string]interface{}{
					"code":    -32601,
					"message": "the method eth_blockNumberMistake does not exist/is not available",
				},
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				t.Errorf("error writing response: %v", err)
			}
			return
		}

		// write the response
		response := `{"jsonrpc":"2.0","id":2,"error":{"code":-32601,"message":"the method eth_blockNumberMistake does not exist/is not available"}}`
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(response))
		if err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	defer server.Close()

	client := server.Client()
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumberMistake",
		"id":      2,
	}
	resp, err := makeRPCRequest(client, server.URL, reqBody)
	if err != nil {
		t.Fatalf("MakeRPCRequest returned unexpected error: %v", err)
	}

	expectedResp := `{"jsonrpc":"2.0","id":2,"error":{"code":-32601,"message":"the method eth_blockNumberMistake does not exist/is not available"}}`
	if string(resp) != expectedResp {
		t.Errorf("expected response %q, got %q", expectedResp, resp)
	}
}

func TestRPC_eth_getBlockByNumber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check the request parameters
		if r.Method != http.MethodPost {
			t.Errorf("expected %s request, got %s", http.MethodPost, r.Method)
		}
		if r.URL.Path != "/" {
			t.Errorf("expected request to path /, got %s", r.URL.Path)
		}

		// write the response
		response := `{"jsonrpc":"2.0","id":2,"result":{"difficulty":"0xd","extraData":"0xd682020983626f7288676f312e31372e32856c696e75780000000000000000004a21925484239cd04672b6d86c9ea2737851a8f87bdb783ee101c830ddcc142018cc108524fb3708a4e7cf82769569b2d8af1e6be85b640d313af8151f2d30c401","gasLimit":"0x1312d00","gasUsed":"0xe13554","hash":"0xe1efb3e3e0e76e7578a6c9216755bf25d22cb0c43dff9aff4f62de507e846d4f","logsBloom":"0x4777a3aad9105f4b34e89c30b567e0e585e96b88cc7d15ca0a94081a3bb2733a2097198a26c8bc19d5bd533cc04d01054984d1a619d3a0a8801215c576793bb93855d9079e32f1ab94f3e9a9bd2f45b299128dee81d530c0783ba55e934c3e19c49a538de200a4c82ce9a8950f1cbce56c844224405d4643952510b411c308a91591ffa4ec1882d214993d2053664f7d49592fb19826758e0e1aadc021483b32eac2b29d10f2cbcd098e4df331e6a501d27c68289a0408f9000ce1385a20e3ded32b0626dece24c200b783a060b29cd048f121000fcbaeb2323ee1b7a48bb14d50579cf9374c54cc22fe9116a141d8369f8df92245ca90f961517b30ee131d0e","miner":"0x0000000000000000000000000000000000000000","mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000","nonce":"0x0000000000000000","number":"0x134e82a","parentHash":"0xa69903bcde35192f34a89e913c67832b88ecc408cf7c376916e32f8c4e9db9a9","receiptsRoot":"0xf54bf69fdc660078853ec0baa2dd78f76b6dd76b1a65fc24dd4eea48c29e5945","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0xf204","stateRoot":"0x01b7c77bb79ec556bfee818ed2bd48ef4f5e0a9a9ff91f1d5a57febb5cf0a6e6","timestamp":"0x61698316","totalDifficulty":"0xe18f426","transactions":["0x50e7d90746c62550262e436912b3e6e7d55cdcfbbfa53d7299f4c68c48ddf050","0xf3dae2f07eb2f695267e543a1b50f0d7e523c7896d536852a59e204312abb75b","0x90c5df7d4a3e4707bda11302a10f91e4b4ed0d995f02e01f6fcae0d719ca90ed","0xb4f9f1be4bc72964d54d14fb3df53697ed099d790ff84344a1ac6755b1bcd68d","0x7355610653e2e5ed3247712afbd00fc8c344420833c297ba339be6e49073b15d","0x24283918af2e7683bfc6dbd5b61d8150fcbb793ae26462d6519a5297444a913f","0xf8ed8714a13b43aa7f22674c3f26e28330a92ccb6f05dfb088116c1d60fe649d","0xbd293b73581896c7ee57f71d19b5edd33becee5d3e86ac937f263cb2f32f4fe1","0x864fc64ba3385380f270435acf647d4df04c8c410c641d1400113ccd18accfb9","0x8408cca9d831b4f449d85b156a240f86e85d4b4821cb88532d7f8b261fefcc25","0x37182f8f7d46accd7248d855846539a029156bf374b20996f0419ea2ff73af80","0xefea4684a9dc8197331c69b85ff0ae05285065f8bd6ef7df50d2e67b7d170b1c","0x835e91250b37f8e09ffbad63148bc5f70d3b3443dfb474fdcd7bbf4961712415","0x21af142d9cef0cdc49da177264dde85288e91073e662e38472d4b6d5095587d0","0x3f012118be7d81dc31820696b737e54776b8965b3ace49b4a4fcb5375113232a","0x3d8de42d30dc80d287a70c7f1be94be3fc6c53501ce6bf54bfd9fffeadc60e82","0x55550023170c8938d9cd9308caf5dfc95d3fd9f667c4afd7bfedddea3c608590","0x4c631c1b1cc5b0bce070ef6e1acddcd125da169f89d1f408987ea50875951e7a","0xccd2ecd7582948c55458f2dc07308ec257659449cc9fd9e3ed511b3e95cb71f6","0xd9d8da80b4c3c24071e5cf4b9b8b179030bcc826c7838e72fa49e43ff5d7ecc0","0x6ac60535c76ed486e06a61839a514182813d39b71c594f23fd167460e37af20a","0x81a9d05298b2a829b61bfbc2887b1c7ad285e9299fe015ce9feb68fba9111bdc","0xf99f1682cbe0c902a1e3206b71110c891752fe902ceb7510243f8eb398be1640","0x00d26ece854448f820c1ec5079d9de1ddba42b82761770e7fb58c537aec31f10","0x67961fd0f666b79e349c1476ba59e3d0c1f7db45b8687ec5c8d6b67c84995e7b","0xb7eaeb3c8a3eaa565563d47d001bb5d7ddafc000d32b85eb551a7e10e2a8ed23","0xe335f258579a6129f8b773809a6da1fe1ac0504a5ed340b7580de9f4c8c598b2","0x3a47fb24c9d582ed401f5a4ea3b1ebbec101c89f445d330946afce8f120d049d","0xad0c546182ca50a73590c3be836032e8942283075fc3c7d5fd2c263219844b64","0x9bb096003a91d94a1b1688760fb8cbc2907060faabee4fcf97ac37e91b1487a9","0xf78400e91d963a6408c46ca799387780bcf84aed29b2eadcf05001d83eb5924e","0x18cf625e8da18e7506cec00863cf6637b41c8e1fc7113b91b34ccf62233a90ca","0x2d4c49ff938cf0fedcda289999fd6a962e6133de15aa635abd1343c06b692f6b","0x3fcc33742ea4187d3a4fd60149346cb07c5676e03a849f6056d098864d828f53","0x8c9cd2910dd3d7682a7b2da858b41a3fbd2b6847aca5c932a292883ae0d173cc","0xf8eaed48880d9bdfd5ee47e0cee9e8edb9eceb99383ef58624571ce858c887aa","0xde7d6d3491d7691a7c7998e61009b69c80331b573d0c516bf73c2a82df647297","0x65d448f151fabae5d6bd64e31d6587e553a026167277b09ee320d2f98e5cdba0","0x704927b379618cc62aeb5d10bda898be9a2b2a3a199bdbac944966aa2b7a7ef3","0x0681389870eb86a7e6d74e98a0ec41cf6e6dd750e063ae6d50d397c89689929e","0x325b07f012fcc6f352cce55084fb58c3e1023d3556f1ea06029d673abd449166","0xaab3b07db1d57e554b59eabe301897c479409a7ef05b65fd851ea2680a2ce5ae","0xd0e10fc6c29b124b0056165ba23339e39b7246889aafcc0e8980fd4700c8fcc2","0x48719f16d2261d896304116b1be4d703017d72053814b376e215246d495d862f","0x51aa809fe2ec3047071925fa8ae8da103123621e192fa3765e7639b3781c31f2","0x9e5a4671e00ee8d1d5274bdac4e6ae390d4adc22588ffbdca4eeb8d6a4d4a988","0x0ae45bf31d2b89fa61be7217c0b888c33e95735c65bd2a9efa1a67cc61576619","0xab219c2d12c379e62f66b59502cca2bde3a4c79d0f0d0e307731a91a009bc146","0x013f459aca7be8a4ce380925fc7aa8acedea8fbfecb6d7bf60c1a457c4fd71c3","0x8a746b59c749307f02b9dda1593a03d14fa85c4cb02853f5b028309b08fe712c","0xf27c978046224fd7bb993b71b5b1b137071c803a6d2837d3f4b4bf1878cdb066","0x765e0863b0755c6c1eb3fad2771af54835765808298188b64c85296d3c59ca8e","0x9b074ba8014b74e6ea6212c4dd9d8f4424fcba58c1213355a5d1d59ffa6bc4ce","0xb6253586d32fd651bb89c55cf706fafb0f64f7d2b3271aff644219691db01ec9","0x684959512d1a4db2bdb0f7e2689e39f853e2b679cfc209549bc87eeb79d057e5","0x8194f82a88a279f6e1269888ebe1d9a2f8baf6913ca85f69898c60c2e463a193","0xa65095224019a6a11e83465024b05bf8ed3d5a4cbf97ce6cc9c5c323cc91366d","0xac91890927c42573ce14242d6fac0df4f8d71af6eb854849c8567b0924c53b19","0xefd26377a4b362b79e5be459431c910225715b76363ebb1a63433d25b7d27e7e","0x8a97a09c5513e5666c6f4ff1c6d742c8cd5cf0132f0a08d0f4ff74fa1df9df54","0xe42bb0c10dc47b8bcc606b5b2040cfb8a0a65b03a33b5e4a1bd5d70b178f4aad","0xb4233efeff30d3b9b1a77e435a75b34cda1e546370b39cf34f6fe46994f087f7","0x9435a9901eec202c070177514b073da96e601e63ab5a23faef21c9364680c7fd","0xb0d6c2d7d8471b2bb22b96953b75582113f190afb7126f691941873fb446aba4","0xe48859a6483bf6a4ffa8d907ca6272741c10a69723200863244df782ee660604","0x528ebd6c3b16391c7301bac1fcf53724d21e7c265a53352ac96b81ddb89f2b60","0xdc45f14e49f99b762dcffa6c8915cae393c8f089f9d90ccfa15631ceddfd1d61"],"transactionsRoot":"0x7c630bf5670f8258a69f0cf6059c674e8def834c370fa58797acee3340563549","uncles":[]}}`
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(response))
		if err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	defer server.Close()

	client := server.Client()
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"id":      2,
	}
	resp, err := makeRPCRequest(client, server.URL, reqBody)
	if err != nil {
		t.Fatalf("MakeRPCRequest returned unexpected error: %v", err)
	}

	expectedResp := `{"jsonrpc":"2.0","id":2,"result":{"difficulty":"0xd","extraData":"0xd682020983626f7288676f312e31372e32856c696e75780000000000000000004a21925484239cd04672b6d86c9ea2737851a8f87bdb783ee101c830ddcc142018cc108524fb3708a4e7cf82769569b2d8af1e6be85b640d313af8151f2d30c401","gasLimit":"0x1312d00","gasUsed":"0xe13554","hash":"0xe1efb3e3e0e76e7578a6c9216755bf25d22cb0c43dff9aff4f62de507e846d4f","logsBloom":"0x4777a3aad9105f4b34e89c30b567e0e585e96b88cc7d15ca0a94081a3bb2733a2097198a26c8bc19d5bd533cc04d01054984d1a619d3a0a8801215c576793bb93855d9079e32f1ab94f3e9a9bd2f45b299128dee81d530c0783ba55e934c3e19c49a538de200a4c82ce9a8950f1cbce56c844224405d4643952510b411c308a91591ffa4ec1882d214993d2053664f7d49592fb19826758e0e1aadc021483b32eac2b29d10f2cbcd098e4df331e6a501d27c68289a0408f9000ce1385a20e3ded32b0626dece24c200b783a060b29cd048f121000fcbaeb2323ee1b7a48bb14d50579cf9374c54cc22fe9116a141d8369f8df92245ca90f961517b30ee131d0e","miner":"0x0000000000000000000000000000000000000000","mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000","nonce":"0x0000000000000000","number":"0x134e82a","parentHash":"0xa69903bcde35192f34a89e913c67832b88ecc408cf7c376916e32f8c4e9db9a9","receiptsRoot":"0xf54bf69fdc660078853ec0baa2dd78f76b6dd76b1a65fc24dd4eea48c29e5945","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0xf204","stateRoot":"0x01b7c77bb79ec556bfee818ed2bd48ef4f5e0a9a9ff91f1d5a57febb5cf0a6e6","timestamp":"0x61698316","totalDifficulty":"0xe18f426","transactions":["0x50e7d90746c62550262e436912b3e6e7d55cdcfbbfa53d7299f4c68c48ddf050","0xf3dae2f07eb2f695267e543a1b50f0d7e523c7896d536852a59e204312abb75b","0x90c5df7d4a3e4707bda11302a10f91e4b4ed0d995f02e01f6fcae0d719ca90ed","0xb4f9f1be4bc72964d54d14fb3df53697ed099d790ff84344a1ac6755b1bcd68d","0x7355610653e2e5ed3247712afbd00fc8c344420833c297ba339be6e49073b15d","0x24283918af2e7683bfc6dbd5b61d8150fcbb793ae26462d6519a5297444a913f","0xf8ed8714a13b43aa7f22674c3f26e28330a92ccb6f05dfb088116c1d60fe649d","0xbd293b73581896c7ee57f71d19b5edd33becee5d3e86ac937f263cb2f32f4fe1","0x864fc64ba3385380f270435acf647d4df04c8c410c641d1400113ccd18accfb9","0x8408cca9d831b4f449d85b156a240f86e85d4b4821cb88532d7f8b261fefcc25","0x37182f8f7d46accd7248d855846539a029156bf374b20996f0419ea2ff73af80","0xefea4684a9dc8197331c69b85ff0ae05285065f8bd6ef7df50d2e67b7d170b1c","0x835e91250b37f8e09ffbad63148bc5f70d3b3443dfb474fdcd7bbf4961712415","0x21af142d9cef0cdc49da177264dde85288e91073e662e38472d4b6d5095587d0","0x3f012118be7d81dc31820696b737e54776b8965b3ace49b4a4fcb5375113232a","0x3d8de42d30dc80d287a70c7f1be94be3fc6c53501ce6bf54bfd9fffeadc60e82","0x55550023170c8938d9cd9308caf5dfc95d3fd9f667c4afd7bfedddea3c608590","0x4c631c1b1cc5b0bce070ef6e1acddcd125da169f89d1f408987ea50875951e7a","0xccd2ecd7582948c55458f2dc07308ec257659449cc9fd9e3ed511b3e95cb71f6","0xd9d8da80b4c3c24071e5cf4b9b8b179030bcc826c7838e72fa49e43ff5d7ecc0","0x6ac60535c76ed486e06a61839a514182813d39b71c594f23fd167460e37af20a","0x81a9d05298b2a829b61bfbc2887b1c7ad285e9299fe015ce9feb68fba9111bdc","0xf99f1682cbe0c902a1e3206b71110c891752fe902ceb7510243f8eb398be1640","0x00d26ece854448f820c1ec5079d9de1ddba42b82761770e7fb58c537aec31f10","0x67961fd0f666b79e349c1476ba59e3d0c1f7db45b8687ec5c8d6b67c84995e7b","0xb7eaeb3c8a3eaa565563d47d001bb5d7ddafc000d32b85eb551a7e10e2a8ed23","0xe335f258579a6129f8b773809a6da1fe1ac0504a5ed340b7580de9f4c8c598b2","0x3a47fb24c9d582ed401f5a4ea3b1ebbec101c89f445d330946afce8f120d049d","0xad0c546182ca50a73590c3be836032e8942283075fc3c7d5fd2c263219844b64","0x9bb096003a91d94a1b1688760fb8cbc2907060faabee4fcf97ac37e91b1487a9","0xf78400e91d963a6408c46ca799387780bcf84aed29b2eadcf05001d83eb5924e","0x18cf625e8da18e7506cec00863cf6637b41c8e1fc7113b91b34ccf62233a90ca","0x2d4c49ff938cf0fedcda289999fd6a962e6133de15aa635abd1343c06b692f6b","0x3fcc33742ea4187d3a4fd60149346cb07c5676e03a849f6056d098864d828f53","0x8c9cd2910dd3d7682a7b2da858b41a3fbd2b6847aca5c932a292883ae0d173cc","0xf8eaed48880d9bdfd5ee47e0cee9e8edb9eceb99383ef58624571ce858c887aa","0xde7d6d3491d7691a7c7998e61009b69c80331b573d0c516bf73c2a82df647297","0x65d448f151fabae5d6bd64e31d6587e553a026167277b09ee320d2f98e5cdba0","0x704927b379618cc62aeb5d10bda898be9a2b2a3a199bdbac944966aa2b7a7ef3","0x0681389870eb86a7e6d74e98a0ec41cf6e6dd750e063ae6d50d397c89689929e","0x325b07f012fcc6f352cce55084fb58c3e1023d3556f1ea06029d673abd449166","0xaab3b07db1d57e554b59eabe301897c479409a7ef05b65fd851ea2680a2ce5ae","0xd0e10fc6c29b124b0056165ba23339e39b7246889aafcc0e8980fd4700c8fcc2","0x48719f16d2261d896304116b1be4d703017d72053814b376e215246d495d862f","0x51aa809fe2ec3047071925fa8ae8da103123621e192fa3765e7639b3781c31f2","0x9e5a4671e00ee8d1d5274bdac4e6ae390d4adc22588ffbdca4eeb8d6a4d4a988","0x0ae45bf31d2b89fa61be7217c0b888c33e95735c65bd2a9efa1a67cc61576619","0xab219c2d12c379e62f66b59502cca2bde3a4c79d0f0d0e307731a91a009bc146","0x013f459aca7be8a4ce380925fc7aa8acedea8fbfecb6d7bf60c1a457c4fd71c3","0x8a746b59c749307f02b9dda1593a03d14fa85c4cb02853f5b028309b08fe712c","0xf27c978046224fd7bb993b71b5b1b137071c803a6d2837d3f4b4bf1878cdb066","0x765e0863b0755c6c1eb3fad2771af54835765808298188b64c85296d3c59ca8e","0x9b074ba8014b74e6ea6212c4dd9d8f4424fcba58c1213355a5d1d59ffa6bc4ce","0xb6253586d32fd651bb89c55cf706fafb0f64f7d2b3271aff644219691db01ec9","0x684959512d1a4db2bdb0f7e2689e39f853e2b679cfc209549bc87eeb79d057e5","0x8194f82a88a279f6e1269888ebe1d9a2f8baf6913ca85f69898c60c2e463a193","0xa65095224019a6a11e83465024b05bf8ed3d5a4cbf97ce6cc9c5c323cc91366d","0xac91890927c42573ce14242d6fac0df4f8d71af6eb854849c8567b0924c53b19","0xefd26377a4b362b79e5be459431c910225715b76363ebb1a63433d25b7d27e7e","0x8a97a09c5513e5666c6f4ff1c6d742c8cd5cf0132f0a08d0f4ff74fa1df9df54","0xe42bb0c10dc47b8bcc606b5b2040cfb8a0a65b03a33b5e4a1bd5d70b178f4aad","0xb4233efeff30d3b9b1a77e435a75b34cda1e546370b39cf34f6fe46994f087f7","0x9435a9901eec202c070177514b073da96e601e63ab5a23faef21c9364680c7fd","0xb0d6c2d7d8471b2bb22b96953b75582113f190afb7126f691941873fb446aba4","0xe48859a6483bf6a4ffa8d907ca6272741c10a69723200863244df782ee660604","0x528ebd6c3b16391c7301bac1fcf53724d21e7c265a53352ac96b81ddb89f2b60","0xdc45f14e49f99b762dcffa6c8915cae393c8f089f9d90ccfa15631ceddfd1d61"],"transactionsRoot":"0x7c630bf5670f8258a69f0cf6059c674e8def834c370fa58797acee3340563549","uncles":[]}}`
	if string(resp) != expectedResp {
		t.Errorf("expected response %q, got %q", expectedResp, resp)
	}
}
