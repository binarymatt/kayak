// @generated by protoc-gen-es v2.3.0 with parameter "target=ts"
// @generated from file kayak/v1/admin.proto (package kayak.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file kayak/v1/admin.proto.
 */
export const file_kayak_v1_admin: GenFile = /*@__PURE__*/
  fileDesc("ChRrYXlhay92MS9hZG1pbi5wcm90bxIIa2F5YWsudjEiPgoPQWRkVm90ZXJSZXF1ZXN0EhIKAmlkGAEgASgJQga6SAPIAQESFwoHYWRkcmVzcxgCIAEoCUIGukgDyAEBIhIKEEFkZFZvdGVyUmVzcG9uc2UiDgoMU3RhdHNSZXF1ZXN0Ik4KCkNvbmZpZ0l0ZW0SEAoIc3VmZnJhZ2UYASABKAkSCgoCaWQYAiABKAkSDwoHYWRkcmVzcxgDIAEoCRIRCglpc19sZWFkZXIYBCABKAgiugEKDVN0YXRzUmVzcG9uc2USDQoFc3RhdGUYASABKAkSIwoFbm9kZXMYAiADKAsyFC5rYXlhay52MS5Db25maWdJdGVtEhQKDGxhc3RfY29udGFjdBgDIAEoCRIxCgVzdGF0cxgEIAMoCzIiLmtheWFrLnYxLlN0YXRzUmVzcG9uc2UuU3RhdHNFbnRyeRosCgpTdGF0c0VudHJ5EgsKA2tleRgBIAEoCRINCgV2YWx1ZRgCIAEoCToCOAEiDwoNTGVhZGVyUmVxdWVzdCItCg5MZWFkZXJSZXNwb25zZRIPCgdhZGRyZXNzGAEgASgJEgoKAmlkGAIgASgJMs4BCgxBZG1pblNlcnZpY2USQwoIQWRkVm90ZXISGS5rYXlhay52MS5BZGRWb3RlclJlcXVlc3QaGi5rYXlhay52MS5BZGRWb3RlclJlc3BvbnNlIgASOgoFU3RhdHMSFi5rYXlhay52MS5TdGF0c1JlcXVlc3QaFy5rYXlhay52MS5TdGF0c1Jlc3BvbnNlIgASPQoGTGVhZGVyEhcua2F5YWsudjEuTGVhZGVyUmVxdWVzdBoYLmtheWFrLnYxLkxlYWRlclJlc3BvbnNlIgBCjQEKDGNvbS5rYXlhay52MUIKQWRtaW5Qcm90b1ABWjBnaXRodWIuY29tL2JpbmFyeW1hdHQva2F5YWsvZ2VuL2theWFrL3YxO2theWFrdjGiAgNLWFiqAghLYXlhay5WMcoCCEtheWFrXFYx4gIUS2F5YWtcVjFcR1BCTWV0YWRhdGHqAglLYXlhazo6VjFiBnByb3RvMw", [file_buf_validate_validate]);

/**
 * @generated from message kayak.v1.AddVoterRequest
 */
export type AddVoterRequest = Message<"kayak.v1.AddVoterRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string address = 2;
   */
  address: string;
};

/**
 * Describes the message kayak.v1.AddVoterRequest.
 * Use `create(AddVoterRequestSchema)` to create a new message.
 */
export const AddVoterRequestSchema: GenMessage<AddVoterRequest> = /*@__PURE__*/
  messageDesc(file_kayak_v1_admin, 0);

/**
 * @generated from message kayak.v1.AddVoterResponse
 */
export type AddVoterResponse = Message<"kayak.v1.AddVoterResponse"> & {
};

/**
 * Describes the message kayak.v1.AddVoterResponse.
 * Use `create(AddVoterResponseSchema)` to create a new message.
 */
export const AddVoterResponseSchema: GenMessage<AddVoterResponse> = /*@__PURE__*/
  messageDesc(file_kayak_v1_admin, 1);

/**
 * @generated from message kayak.v1.StatsRequest
 */
export type StatsRequest = Message<"kayak.v1.StatsRequest"> & {
};

/**
 * Describes the message kayak.v1.StatsRequest.
 * Use `create(StatsRequestSchema)` to create a new message.
 */
export const StatsRequestSchema: GenMessage<StatsRequest> = /*@__PURE__*/
  messageDesc(file_kayak_v1_admin, 2);

/**
 * @generated from message kayak.v1.ConfigItem
 */
export type ConfigItem = Message<"kayak.v1.ConfigItem"> & {
  /**
   * @generated from field: string suffrage = 1;
   */
  suffrage: string;

  /**
   * @generated from field: string id = 2;
   */
  id: string;

  /**
   * @generated from field: string address = 3;
   */
  address: string;

  /**
   * @generated from field: bool is_leader = 4;
   */
  isLeader: boolean;
};

/**
 * Describes the message kayak.v1.ConfigItem.
 * Use `create(ConfigItemSchema)` to create a new message.
 */
export const ConfigItemSchema: GenMessage<ConfigItem> = /*@__PURE__*/
  messageDesc(file_kayak_v1_admin, 3);

/**
 * @generated from message kayak.v1.StatsResponse
 */
export type StatsResponse = Message<"kayak.v1.StatsResponse"> & {
  /**
   * @generated from field: string state = 1;
   */
  state: string;

  /**
   * @generated from field: repeated kayak.v1.ConfigItem nodes = 2;
   */
  nodes: ConfigItem[];

  /**
   * @generated from field: string last_contact = 3;
   */
  lastContact: string;

  /**
   * @generated from field: map<string, string> stats = 4;
   */
  stats: { [key: string]: string };
};

/**
 * Describes the message kayak.v1.StatsResponse.
 * Use `create(StatsResponseSchema)` to create a new message.
 */
export const StatsResponseSchema: GenMessage<StatsResponse> = /*@__PURE__*/
  messageDesc(file_kayak_v1_admin, 4);

/**
 * @generated from message kayak.v1.LeaderRequest
 */
export type LeaderRequest = Message<"kayak.v1.LeaderRequest"> & {
};

/**
 * Describes the message kayak.v1.LeaderRequest.
 * Use `create(LeaderRequestSchema)` to create a new message.
 */
export const LeaderRequestSchema: GenMessage<LeaderRequest> = /*@__PURE__*/
  messageDesc(file_kayak_v1_admin, 5);

/**
 * @generated from message kayak.v1.LeaderResponse
 */
export type LeaderResponse = Message<"kayak.v1.LeaderResponse"> & {
  /**
   * @generated from field: string address = 1;
   */
  address: string;

  /**
   * @generated from field: string id = 2;
   */
  id: string;
};

/**
 * Describes the message kayak.v1.LeaderResponse.
 * Use `create(LeaderResponseSchema)` to create a new message.
 */
export const LeaderResponseSchema: GenMessage<LeaderResponse> = /*@__PURE__*/
  messageDesc(file_kayak_v1_admin, 6);

/**
 * @generated from service kayak.v1.AdminService
 */
export const AdminService: GenService<{
  /**
   * @generated from rpc kayak.v1.AdminService.AddVoter
   */
  addVoter: {
    methodKind: "unary";
    input: typeof AddVoterRequestSchema;
    output: typeof AddVoterResponseSchema;
  },
  /**
   * @generated from rpc kayak.v1.AdminService.Stats
   */
  stats: {
    methodKind: "unary";
    input: typeof StatsRequestSchema;
    output: typeof StatsResponseSchema;
  },
  /**
   * @generated from rpc kayak.v1.AdminService.Leader
   */
  leader: {
    methodKind: "unary";
    input: typeof LeaderRequestSchema;
    output: typeof LeaderResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_kayak_v1_admin, 0);

