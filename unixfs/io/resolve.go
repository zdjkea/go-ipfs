package io

import (
	"context"

	dag "github.com/ipfs/go-ipfs/merkledag"
	ft "github.com/ipfs/go-ipfs/unixfs"
	hamt "github.com/ipfs/go-ipfs/unixfs/hamt"

	node "gx/ipfs/QmUsVJ7AEnGyjX8YWnrwq9vmECVGwBQNAKPpgz5KSg8dcq/go-ipld-node"
)

func ResolveUnixfsOnce(ctx context.Context, ds dag.DAGService, nd node.Node, name string) (*node.Link, error) {
	switch nd := nd.(type) {
	case *dag.ProtoNode:
		upb, err := ft.FromBytes(nd.Data())
		if err != nil {
			// Not a unixfs node, use standard object traversal code
			return nd.GetNodeLink(name)
		}

		switch upb.GetType() {
		case ft.THAMTShard:
			s, err := hamt.NewHamtFromDag(ds, nd)
			if err != nil {
				return nil, err
			}

			// TODO: optimized routine on HAMT for returning a dag.Link to avoid extra disk hits
			out, err := s.Find(ctx, name)
			if err != nil {
				return nil, err
			}

			return node.MakeLink(out)
		default:
			return nd.GetNodeLink(name)
		}
	default:
		lnk, _, err := nd.ResolveLink([]string{name})
		if err != nil {
			return nil, err
		}
		return lnk, nil
	}
}
