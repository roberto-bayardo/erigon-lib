package types

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/ledgerwatch/erigon-lib/common/hexutility"
)

var (
	// blobTxNetworkWrapperHex is the network-encoding (ssz) of the following tx, including 2 blobs
	// (not displayed in the json below) that should pass verification.
	//
	// {
	//   "nonce": "0xa",
	//   "gasPrice": null,
	//   "maxPriorityFeePerGas": "0x2a",
	//   "maxFeePerGas": "0xa",
	//   "gas": "0x1e241",
	//   "value": "0x64",
	//   "input": "0x616263646566",
	//   "v": "0x1",
	//   "r": "0xe995f26f2f424703e00ef9c9709248dc6587f3045e2dd536eedf96651a4b680d",
	//   "s": "0x13836dded49612eb61c61e9c61aa343a26f4ba37b5d53da3b6d9326b64a09668",
	//   "to": "0x095e7baea6a6c7c4c2dfeb977efac326af552d87",
	//   "chainId": "0x1",
	//   "accessList": [
	//     {
	//       "address": "0x0000000000000000000000000000000000000001",
	//       "storageKeys": [
	//         "0x0000000000000000000000000000000000000000000000000000000000000000",
	//         "0x0100000000000000000000000000000000000000000000000000000000000000"
	//       ]
	//     },
	//     {
	//       "address": "0x0000000000000000000000000000000000000002",
	//       "storageKeys": [
	//         "0x0200000000000000000000000000000000000000000000000000000000000000"
	//       ]
	//     }
	//   ],
	//   "maxFeePerDataGas": "0x0",
	//   "blobVersionedHashes": [ ... ],
	//   "kzgAggregatedProof": "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	//   "hash": "0xabfab29ef05293b52e448f5e85eae4a99c1496cdf59f59987a37ba90912c8801"
	// }
	//
	//go:embed testdata/blobtx.txt
	blobTxHex string

	// Sender address of the blobTxHex transaction
	blobTxSender string = "a614f00e912255dae37395ee7ba37a095dc98f81"

	// blobTxCreateHex is just like the previous test case only the To: field is nil (simulating
	// contract creation) and there are no blobs
	blobTxCreateHex string = "10000000b4010000b4010000b4010000450000000137d13cd1b1e351f12556ef6a6fa74142f5cec279aa50df41728757f2c1203e6e67db95d38790467fe12da748c6e1d7f22319649b94ae64d3fa2e4d5d5ae6d25a01000000000000000000000000000000000000000000000000000000000000000a000000000000002a000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000041e2010000000000c00000006400000000000000000000000000000000000000000000000000000000000000c1000000c700000000000000000000000000000000000000000000000000000000000000000000005f010000006162636465660800000060000000000000000000000000000000000000000000000118000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002180000000200000000000000000000000000000000000000000000000000000000000000"
)

func TestDeserializeBlobTx(t *testing.T) {
	txBytes := hexutility.MustDecodeHex(strings.TrimSpace(blobTxHex))
	w := wrapper{}
	err := w.Deserialize(txBytes)
	if err != nil {
		t.Fatalf("couldn't deserialize blob tx: %v", err)
	}

	if w.chainID.Uint64() != 1 {
		t.Errorf("Expected chain id == 1, got: %v", w.chainID.Uint64())
	}
	if w.nonce != 10 {
		t.Errorf("Expected nonce == 10, got: %v", w.nonce)
	}
	if w.maxPriorityFeePerGas.Uint64() != 42 {
		t.Errorf("Expected maxPriorityFeePerGas == 42, got %v", w.maxPriorityFeePerGas.Uint64())
	}
	if w.maxFeePerGas.Uint64() != 10 {
		t.Errorf("Expected maxFeePerGas == 10, got %v", w.maxFeePerGas.Uint64())
	}
	if w.gas != 123457 {
		t.Errorf("Expected gas == 123457, got %v", w.gas)
	}
	if w.creation == true {
		t.Errorf("Expected !w.creation")
	}
	if w.value.Uint64() != 100 {
		t.Errorf("Expected w.value == 100, got %v", w.value.Uint64())
	}
	if w.dataLen != 6 {
		t.Errorf("Expected dataLen == 6, got %v", w.dataLen)
	}
	if w.dataNonZeroLen != 6 {
		t.Errorf("Expected dataNonZeroLen == 6, got %v", w.dataNonZeroLen)
	}
	if w.numBlobHashes != 2 {
		t.Errorf("Expected 2 blob hashes, got %v", w.numBlobHashes)
	}
	if w.accessListAddressCount != 2 {
		t.Errorf("Expected 2 addresses in access list, got %v", w.accessListAddressCount)
	}
	if w.accessListKeyCount != 3 {
		t.Errorf("Expected 3 keys in access list, got %v", w.accessListKeyCount)
	}

	/// TODO: broken by crate-crypto / spec updates
	//err = w.VerifyBlobs(txBytes)
	//if err != nil {
	//t.Errorf("blob verification failed: %v", err)
	//}

	// Now mangle a proof byte and make sure verification fails
	// TODO: FIX these tests with updated blob tx spec
	// oldByte := txBytes[w.proofsOffset]
	// txBytes[w.proofsOffset+1] = 0xff
	// err = w.VerifyBlobs(txBytes)
	// if err == nil {
	// 	t.Errorf("expected blob verification to fail")
	// }
	// t.Logf("Got error as expected: %v", err)
	// txBytes[w.proofsOffset+1] = oldByte // restore the mangled byte

	// // Now mangle a blob byte and make sure verification fails
	// txBytes[w.blobsOffset] = 0xff
	// err = w.VerifyBlobs(txBytes)
	// if err == nil {
	// 	t.Errorf("expected blob verification to fail")
	// }
	// t.Logf("Got error as expected: %v", err)
}

func TestDeserializeBlobCreateTx(t *testing.T) {
	txBytes := hexutility.MustDecodeHex(strings.TrimSpace(blobTxCreateHex))
	w := wrapper{}
	err := w.Deserialize(txBytes)
	if err != nil {
		t.Fatalf("couldn't deserialize blob tx: %v", err)
	}

	// Test the different outcomes from the previous test case
	if !w.creation {
		t.Errorf("Expected w.creation")
	}
	if w.numBlobHashes != 0 {
		t.Errorf("Expected 0 blob hashes, got %v", w.numBlobHashes)
	}

	// This should fail since the tx has 0 blobs
	err = w.VerifyBlobs(txBytes)
	if err == nil {
		t.Errorf("blob verification should have failed due to 0 blobs")
	}

	// Remaining fields should be same as last test case
	if w.chainID.Uint64() != 1 {
		t.Errorf("Expected chain id == 1, got: %v", w.chainID.Uint64())
	}
	if w.nonce != 10 {
		t.Errorf("Expected nonce == 10, got: %v", w.nonce)
	}
	if w.maxPriorityFeePerGas.Uint64() != 42 {
		t.Errorf("Expected maxPriorityFeePerGas == 42, got %v", w.maxPriorityFeePerGas.Uint64())
	}
	if w.maxFeePerGas.Uint64() != 10 {
		t.Errorf("Expected maxFeePerGas == 10, got %v", w.maxFeePerGas.Uint64())
	}
	if w.gas != 123457 {
		t.Errorf("Expected gas == 123457, got %v", w.gas)
	}
	if w.value.Uint64() != 100 {
		t.Errorf("Expected w.value == 100, got %v", w.value.Uint64())
	}
	if w.dataLen != 6 {
		t.Errorf("Expected dataLen == 6, got %v", w.dataLen)
	}
	if w.dataNonZeroLen != 6 {
		t.Errorf("Expected dataNonZeroLen == 6, got %v", w.dataNonZeroLen)
	}
	if w.accessListAddressCount != 2 {
		t.Errorf("Expected 2 addresses in access list, got %v", w.accessListAddressCount)
	}
	if w.accessListKeyCount != 3 {
		t.Errorf("Expected 3 keys in access list, got %v", w.accessListKeyCount)
	}
}
