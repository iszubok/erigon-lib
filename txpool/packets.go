/*
   Copyright 2021 Erigon contributors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package txpool

import (
	"fmt"

	"github.com/ledgerwatch/erigon-lib/rlp"
)

type NewPooledTransactionHashesPacket [][32]byte

// ParseHashesCount looks at the RLP length Prefix for list of 32-byte hashes
// and returns number of hashes in the list to expect
func ParseHashesCount(payload Hashes, pos int) (int, int, error) {
	dataPos, dataLen, err := rlp.List(payload, pos)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: hashes len: %w", rlp.ParseHashErrorPrefix, err)
	}
	if dataLen%33 != 0 {
		return 0, 0, fmt.Errorf("%s: hashes len must be multiple of 33", rlp.ParseHashErrorPrefix)
	}
	return dataLen / 33, dataPos, nil
}

// EncodeHashes produces RLP encoding of given number of hashes, as RLP list
// It appends encoding to the given given slice (encodeBuf), reusing the space
// there is there is enough capacity.
// The first returned value is the slice where encodinfg
func EncodeHashes(hashes []byte, encodeBuf []byte) []byte {
	hashesLen := len(hashes) / 32 * 33
	dataLen := hashesLen
	encodeBuf = ensureEnoughSize(encodeBuf, rlp.ListPrefixLen(hashesLen)+dataLen)
	rlp.EncodeHashes(hashes, encodeBuf)
	return encodeBuf
}

func ensureEnoughSize(in []byte, size int) []byte {
	if cap(in) < size {
		return make([]byte, size)
	}
	return in[:size] // Reuse the space in pkbuf, is it has enough capacity
}

// EncodeGetPooledTransactions66 produces encoding of GetPooledTransactions66 packet
func EncodeGetPooledTransactions66(hashes []byte, requestId uint64, encodeBuf []byte) ([]byte, error) {
	pos := 0
	hashesLen := len(hashes) / 32 * 33
	dataLen := rlp.ListPrefixLen(hashesLen) + hashesLen + rlp.U64Len(requestId)
	encodeBuf = ensureEnoughSize(encodeBuf, rlp.ListPrefixLen(dataLen)+dataLen)
	// Length Prefix for the entire structure
	pos += rlp.EncodeListPrefix(dataLen, encodeBuf[pos:])
	pos += rlp.EncodeU64(requestId, encodeBuf[pos:])
	pos += rlp.EncodeHashes(hashes, encodeBuf[pos:])
	return encodeBuf, nil
}