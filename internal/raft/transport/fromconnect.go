package transport

import (
	"github.com/hashicorp/raft"

	transportv1 "github.com/binarymatt/kayak/gen/transport/v1"
)

func decodeRPCHeader(resp *transportv1.RPCHeader) raft.RPCHeader {
	return raft.RPCHeader{
		ProtocolVersion: raft.ProtocolVersion(resp.ProtocolVersion),
		ID:              resp.Id,
		Addr:            resp.Addr,
	}
}
func decodeInstallSnapshotRequest(m *transportv1.InstallSnapshotRequest) *raft.InstallSnapshotRequest {
	return &raft.InstallSnapshotRequest{
		RPCHeader:          decodeRPCHeader(m.RpcHeader),
		SnapshotVersion:    raft.SnapshotVersion(m.SnapshotVersion),
		Term:               m.Term,
		Leader:             m.Leader,
		LastLogIndex:       m.LastLogIndex,
		LastLogTerm:        m.LastLogTerm,
		Peers:              m.Peers,
		Configuration:      m.Configuration,
		ConfigurationIndex: m.ConfigurationIndex,
		Size:               m.Size,
	}
}
func decodeInstallSnapshotResponse(resp *transportv1.InstallSnapshotResponse) *raft.InstallSnapshotResponse {
	return &raft.InstallSnapshotResponse{
		RPCHeader: decodeRPCHeader(resp.RpcHeader),
		Term:      resp.Term,
		Success:   resp.Success,
	}
}
func decodeRequestPreVoteRequest(m *transportv1.RequestPreVoteRequest) *raft.RequestPreVoteRequest {
	return &raft.RequestPreVoteRequest{
		RPCHeader:    decodeRPCHeader(m.RpcHeader),
		Term:         m.Term,
		LastLogIndex: m.LastLogIndex,
		LastLogTerm:  m.LastLogTerm,
	}
}

func decodeRequestPreVoteResponse(resp *transportv1.RequestPreVoteResponse) *raft.RequestPreVoteResponse {
	return &raft.RequestPreVoteResponse{
		RPCHeader: decodeRPCHeader(resp.RpcHeader),
		Term:      resp.Term,
		Granted:   resp.Granted,
	}
}
func decodeRequestVoteRequest(m *transportv1.RequestVoteRequest) *raft.RequestVoteRequest {
	return &raft.RequestVoteRequest{
		RPCHeader:          decodeRPCHeader(m.RpcHeader),
		Term:               m.Term,
		Candidate:          m.Candidate,
		LastLogIndex:       m.LastLogIndex,
		LastLogTerm:        m.LastLogTerm,
		LeadershipTransfer: m.LeadershipTransfer,
	}
}
func decodeRequestVoteResponse(resp *transportv1.RequestVoteResponse) *raft.RequestVoteResponse {
	return &raft.RequestVoteResponse{
		RPCHeader: decodeRPCHeader(resp.RpcHeader),
		Term:      resp.Term,
		Peers:     resp.Peers,
		Granted:   resp.Granted,
	}
}
func decodeTimeoutNowRequest(m *transportv1.TimeoutNowRequest) *raft.TimeoutNowRequest {
	return &raft.TimeoutNowRequest{
		RPCHeader: decodeRPCHeader(m.RpcHeader),
	}
}
func decodeTimeoutNowResponse(resp *transportv1.TimeoutNowResponse) *raft.TimeoutNowResponse {
	return &raft.TimeoutNowResponse{
		RPCHeader: decodeRPCHeader(resp.RpcHeader),
	}
}

func decodeAppendEntriesRequest(req *transportv1.AppendEntriesRequest) *raft.AppendEntriesRequest {
	return &raft.AppendEntriesRequest{
		RPCHeader:         decodeRPCHeader(req.RpcHeader),
		Term:              req.Term,
		Leader:            req.Leader,
		PrevLogEntry:      req.PrevLogEntry,
		PrevLogTerm:       req.PrevLogTerm,
		Entries:           decodeLogs(req.Entries),
		LeaderCommitIndex: req.LeaderCommitIndex,
	}
}

func decodeLogs(l []*transportv1.Log) []*raft.Log {
	ret := make([]*raft.Log, len(l))
	for i, l := range l {
		ret[i] = decodeLog(l)
	}
	return ret
}

func decodeLog(l *transportv1.Log) *raft.Log {
	return &raft.Log{
		Index:      l.Index,
		Term:       l.Term,
		Type:       decodeLogType(l.Type),
		Data:       l.Data,
		Extensions: l.Extensions,
		AppendedAt: l.AppendedAt.AsTime(),
	}
}

func decodeLogType(m transportv1.Log_LogType) raft.LogType {
	switch m {
	case transportv1.Log_LOG_COMMAND:
		return raft.LogCommand
	case transportv1.Log_LOG_NOOP:
		return raft.LogNoop
	case transportv1.Log_LOG_ADD_PEER_DEPRECATED:
		return raft.LogAddPeerDeprecated
	case transportv1.Log_LOG_REMOVE_PEER_DEPRECATED:
		return raft.LogRemovePeerDeprecated
	case transportv1.Log_LOG_BARRIER:
		return raft.LogBarrier
	case transportv1.Log_LOG_CONFIGURATION:
		return raft.LogConfiguration
	default:
		panic("invalid LogType")
	}
}

func decodeAppendEntriesResponse(resp *transportv1.AppendEntriesResponse) *raft.AppendEntriesResponse {
	return &raft.AppendEntriesResponse{
		RPCHeader:      decodeRPCHeader(resp.RpcHeader),
		Term:           resp.Term,
		LastLog:        resp.LastLog,
		Success:        resp.Success,
		NoRetryBackoff: resp.NoRetryBackoff,
	}
}
