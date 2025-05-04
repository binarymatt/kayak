package transport

import (
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/types/known/timestamppb"

	transportv1 "github.com/binarymatt/kayak/gen/transport/v1"
)

func encodeAppendEntriesRequest(s *raft.AppendEntriesRequest) *transportv1.AppendEntriesRequest {
	return &transportv1.AppendEntriesRequest{
		RpcHeader:         encodeRPCHeader(s.RPCHeader),
		Term:              s.Term,
		Leader:            s.Addr,
		PrevLogEntry:      s.PrevLogEntry,
		PrevLogTerm:       s.PrevLogTerm,
		Entries:           encodeLogs(s.Entries),
		LeaderCommitIndex: s.LeaderCommitIndex,
	}
}

func encodeRPCHeader(s raft.RPCHeader) *transportv1.RPCHeader {
	return &transportv1.RPCHeader{
		ProtocolVersion: int64(s.ProtocolVersion),
		Id:              s.ID,
		Addr:            s.Addr,
	}
}

func encodeLogs(s []*raft.Log) []*transportv1.Log {
	ret := make([]*transportv1.Log, len(s))
	for i, l := range s {
		ret[i] = encodeLog(l)
	}
	return ret
}

func encodeLog(s *raft.Log) *transportv1.Log {
	return &transportv1.Log{
		Index:      s.Index,
		Term:       s.Term,
		Type:       encodeLogType(s.Type),
		Data:       s.Data,
		Extensions: s.Extensions,
		AppendedAt: timestamppb.New(s.AppendedAt),
	}
}

func encodeLogType(s raft.LogType) transportv1.Log_LogType {
	switch s {
	case raft.LogCommand:
		return transportv1.Log_LOG_COMMAND
	case raft.LogNoop:
		return transportv1.Log_LOG_NOOP
	case raft.LogAddPeerDeprecated:
		return transportv1.Log_LOG_ADD_PEER_DEPRECATED
	case raft.LogRemovePeerDeprecated:
		return transportv1.Log_LOG_REMOVE_PEER_DEPRECATED
	case raft.LogBarrier:
		return transportv1.Log_LOG_BARRIER
	case raft.LogConfiguration:
		return transportv1.Log_LOG_CONFIGURATION
	default:
		panic("invalid LogType")
	}
}
func encodeAppendEntriesResponse(s *raft.AppendEntriesResponse) *transportv1.AppendEntriesResponse {
	return &transportv1.AppendEntriesResponse{
		RpcHeader:      encodeRPCHeader(s.RPCHeader),
		Term:           s.Term,
		LastLog:        s.LastLog,
		Success:        s.Success,
		NoRetryBackoff: s.NoRetryBackoff,
	}
}

func encodeInstallSnapshotRequest(req *raft.InstallSnapshotRequest) *transportv1.InstallSnapshotRequest {
	return &transportv1.InstallSnapshotRequest{
		RpcHeader:          encodeRPCHeader(req.RPCHeader),
		SnapshotVersion:    int64(req.SnapshotVersion),
		Term:               req.Term,
		Leader:             req.Leader,
		LastLogIndex:       req.LastLogIndex,
		LastLogTerm:        req.LastLogTerm,
		Peers:              req.Peers,
		Configuration:      req.Configuration,
		ConfigurationIndex: req.ConfigurationIndex,
		Size:               req.Size,
	}
}

func encodeInstallSnapshotResponse(s *raft.InstallSnapshotResponse) *transportv1.InstallSnapshotResponse {
	return &transportv1.InstallSnapshotResponse{
		RpcHeader: encodeRPCHeader(s.RPCHeader),
		Term:      s.Term,
		Success:   s.Success,
	}
}

func encodeRequestPreVoteRequest(s *raft.RequestPreVoteRequest) *transportv1.RequestPreVoteRequest {
	return &transportv1.RequestPreVoteRequest{
		RpcHeader:    encodeRPCHeader(s.RPCHeader),
		Term:         s.Term,
		LastLogIndex: s.LastLogIndex,
		LastLogTerm:  s.LastLogTerm,
	}
}

func encodeRequestPreVoteResponse(s *raft.RequestPreVoteResponse) *transportv1.RequestPreVoteResponse {
	return &transportv1.RequestPreVoteResponse{
		RpcHeader: encodeRPCHeader(s.RPCHeader),
		Term:      s.Term,
		Granted:   s.Granted,
	}
}

func encodeRequestVoteRequest(s *raft.RequestVoteRequest) *transportv1.RequestVoteRequest {
	return &transportv1.RequestVoteRequest{
		RpcHeader:          encodeRPCHeader(s.RPCHeader),
		Term:               s.Term,
		Candidate:          s.Addr,
		LastLogIndex:       s.LastLogIndex,
		LastLogTerm:        s.LastLogTerm,
		LeadershipTransfer: s.LeadershipTransfer,
	}
}
func encodeRequestVoteResponse(s *raft.RequestVoteResponse) *transportv1.RequestVoteResponse {
	return &transportv1.RequestVoteResponse{
		RpcHeader: encodeRPCHeader(s.RPCHeader),
		Term:      s.Term,
		Peers:     s.Peers,
		Granted:   s.Granted,
	}
}

func encodeTimeoutNowRequest(s *raft.TimeoutNowRequest) *transportv1.TimeoutNowRequest {
	return &transportv1.TimeoutNowRequest{
		RpcHeader: encodeRPCHeader(s.RPCHeader),
	}
}

func encodeTimeoutNowResponse(s *raft.TimeoutNowResponse) *transportv1.TimeoutNowResponse {
	return &transportv1.TimeoutNowResponse{
		RpcHeader: encodeRPCHeader(s.RPCHeader),
	}
}
