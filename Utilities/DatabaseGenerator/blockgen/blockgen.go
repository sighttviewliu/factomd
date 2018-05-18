package blockgen

import (
	"github.com/FactomProject/factomd/common/directoryBlock"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/state"
)

type BlockGen struct {
	EntryGenerator  IFullEntryGenerator
	AuthoritySigner IAuthSigner
}

func NewBlockGen(config EntryGeneratorConfig) (*BlockGen, error) {
	b := new(BlockGen)
	b.AuthoritySigner = new(DefaultAuthSigner)

	fkey, err := primitives.NewPrivateKeyFromHex("FB3B471B1DCDADFEB856BD0B02D8BF49ACE0EDD372A3D9F2A95B78EC12A324D6")
	if err != nil {
		return nil, err
	}

	b.EntryGenerator = NewFullEntryGenerator(*primitives.RandomPrivateKey(), *fkey, config)
	return b, nil
}

func (bg *BlockGen) NewBlock(prev *state.DBState, netid uint32) (*state.DBState, error) {
	// ABlock
	nab := bg.AuthoritySigner.SignBlock(prev)
	next := primitives.Timestamp(prev.DirectoryBlock.GetHeader().GetTimestamp().GetTimeMilliUInt64() + 10*60)
	if prev.DirectoryBlock.GetDatabaseHeight() == 0 {
		next = *primitives.NewTimestampNow()
	}

	// Entries (need entries for ecblock)
	newDBState, err := bg.EntryGenerator.NewBlockSet(prev, &next)
	if err != nil {
		return nil, err
	}
	newDBState.AdminBlock = nab
	newDBState.ABHash = nab.DatabasePrimaryIndex()

	// DBlock
	dblock := directoryBlock.NewDirectoryBlock(prev.DirectoryBlock)
	dblock.GetHeader().SetNetworkID(netid)
	dblock.GetHeader().SetTimestamp(&next)
	dblock.SetABlockHash(nab)
	dblock.SetECBlockHash(newDBState.EntryCreditBlock)
	dblock.SetFBlockHash(newDBState.FactoidBlock)

	for _, eb := range newDBState.EntryBlocks {
		k, _ := eb.KeyMR()
		dblock.AddEntry(eb.GetChainID(), k)
	}

	dblock.HeaderHash()
	dblock.BuildBodyMR()
	dblock.BuildKeyMerkleRoot()

	newDBState.DirectoryBlock = dblock
	return newDBState, nil
}

func newDblock(prev interfaces.IDirectoryBlock) interfaces.IDirectoryBlock {
	dblock := directoryBlock.NewDirectoryBlock(prev)
	return dblock
}
