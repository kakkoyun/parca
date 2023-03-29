// @generated by protobuf-ts 2.8.3 with parameter long_type_string,generate_dependencies
// @generated from protobuf file "parca/query/v1alpha1/query.proto" (package "parca.query.v1alpha1", syntax proto3)
// tslint:disable
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { QueryService } from "./query";
import type { ShareProfileResponse } from "./query";
import type { ShareProfileRequest } from "./query";
import type { ValuesResponse } from "./query";
import type { ValuesRequest } from "./query";
import type { LabelsResponse } from "./query";
import type { LabelsRequest } from "./query";
import type { ProfileTypesResponse } from "./query";
import type { ProfileTypesRequest } from "./query";
import type { SeriesResponse } from "./query";
import type { SeriesRequest } from "./query";
import type { QueryResponse } from "./query";
import type { QueryRequest } from "./query";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { QueryRangeResponse } from "./query";
import type { QueryRangeRequest } from "./query";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * QueryService is the service that provides APIs to retrieve and inspect profiles
 *
 * @generated from protobuf service parca.query.v1alpha1.QueryService
 */
export interface IQueryServiceClient {
    /**
     * QueryRange performs a profile query over a time range
     *
     * @generated from protobuf rpc: QueryRange(parca.query.v1alpha1.QueryRangeRequest) returns (parca.query.v1alpha1.QueryRangeResponse);
     */
    queryRange(input: QueryRangeRequest, options?: RpcOptions): UnaryCall<QueryRangeRequest, QueryRangeResponse>;
    /**
     * Query performs a profile query
     *
     * @generated from protobuf rpc: Query(parca.query.v1alpha1.QueryRequest) returns (parca.query.v1alpha1.QueryResponse);
     */
    query(input: QueryRequest, options?: RpcOptions): UnaryCall<QueryRequest, QueryResponse>;
    /**
     * Series is unimplemented
     *
     * @generated from protobuf rpc: Series(parca.query.v1alpha1.SeriesRequest) returns (parca.query.v1alpha1.SeriesResponse);
     */
    series(input: SeriesRequest, options?: RpcOptions): UnaryCall<SeriesRequest, SeriesResponse>;
    /**
     * ProfileTypes returns the list of available profile types.
     *
     * @generated from protobuf rpc: ProfileTypes(parca.query.v1alpha1.ProfileTypesRequest) returns (parca.query.v1alpha1.ProfileTypesResponse);
     */
    profileTypes(input: ProfileTypesRequest, options?: RpcOptions): UnaryCall<ProfileTypesRequest, ProfileTypesResponse>;
    /**
     * Labels returns the set of label names against a given matching string and time frame
     *
     * @generated from protobuf rpc: Labels(parca.query.v1alpha1.LabelsRequest) returns (parca.query.v1alpha1.LabelsResponse);
     */
    labels(input: LabelsRequest, options?: RpcOptions): UnaryCall<LabelsRequest, LabelsResponse>;
    /**
     * Values returns the set of values that match a given label and time frame
     *
     * @generated from protobuf rpc: Values(parca.query.v1alpha1.ValuesRequest) returns (parca.query.v1alpha1.ValuesResponse);
     */
    values(input: ValuesRequest, options?: RpcOptions): UnaryCall<ValuesRequest, ValuesResponse>;
    /**
     * ShareProfile uploads the given profile to pprof.me and returns a link to the profile.
     *
     * @generated from protobuf rpc: ShareProfile(parca.query.v1alpha1.ShareProfileRequest) returns (parca.query.v1alpha1.ShareProfileResponse);
     */
    shareProfile(input: ShareProfileRequest, options?: RpcOptions): UnaryCall<ShareProfileRequest, ShareProfileResponse>;
}
/**
 * QueryService is the service that provides APIs to retrieve and inspect profiles
 *
 * @generated from protobuf service parca.query.v1alpha1.QueryService
 */
export class QueryServiceClient implements IQueryServiceClient, ServiceInfo {
    typeName = QueryService.typeName;
    methods = QueryService.methods;
    options = QueryService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * QueryRange performs a profile query over a time range
     *
     * @generated from protobuf rpc: QueryRange(parca.query.v1alpha1.QueryRangeRequest) returns (parca.query.v1alpha1.QueryRangeResponse);
     */
    queryRange(input: QueryRangeRequest, options?: RpcOptions): UnaryCall<QueryRangeRequest, QueryRangeResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<QueryRangeRequest, QueryRangeResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * Query performs a profile query
     *
     * @generated from protobuf rpc: Query(parca.query.v1alpha1.QueryRequest) returns (parca.query.v1alpha1.QueryResponse);
     */
    query(input: QueryRequest, options?: RpcOptions): UnaryCall<QueryRequest, QueryResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<QueryRequest, QueryResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * Series is unimplemented
     *
     * @generated from protobuf rpc: Series(parca.query.v1alpha1.SeriesRequest) returns (parca.query.v1alpha1.SeriesResponse);
     */
    series(input: SeriesRequest, options?: RpcOptions): UnaryCall<SeriesRequest, SeriesResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<SeriesRequest, SeriesResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * ProfileTypes returns the list of available profile types.
     *
     * @generated from protobuf rpc: ProfileTypes(parca.query.v1alpha1.ProfileTypesRequest) returns (parca.query.v1alpha1.ProfileTypesResponse);
     */
    profileTypes(input: ProfileTypesRequest, options?: RpcOptions): UnaryCall<ProfileTypesRequest, ProfileTypesResponse> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<ProfileTypesRequest, ProfileTypesResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * Labels returns the set of label names against a given matching string and time frame
     *
     * @generated from protobuf rpc: Labels(parca.query.v1alpha1.LabelsRequest) returns (parca.query.v1alpha1.LabelsResponse);
     */
    labels(input: LabelsRequest, options?: RpcOptions): UnaryCall<LabelsRequest, LabelsResponse> {
        const method = this.methods[4], opt = this._transport.mergeOptions(options);
        return stackIntercept<LabelsRequest, LabelsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * Values returns the set of values that match a given label and time frame
     *
     * @generated from protobuf rpc: Values(parca.query.v1alpha1.ValuesRequest) returns (parca.query.v1alpha1.ValuesResponse);
     */
    values(input: ValuesRequest, options?: RpcOptions): UnaryCall<ValuesRequest, ValuesResponse> {
        const method = this.methods[5], opt = this._transport.mergeOptions(options);
        return stackIntercept<ValuesRequest, ValuesResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * ShareProfile uploads the given profile to pprof.me and returns a link to the profile.
     *
     * @generated from protobuf rpc: ShareProfile(parca.query.v1alpha1.ShareProfileRequest) returns (parca.query.v1alpha1.ShareProfileResponse);
     */
    shareProfile(input: ShareProfileRequest, options?: RpcOptions): UnaryCall<ShareProfileRequest, ShareProfileResponse> {
        const method = this.methods[6], opt = this._transport.mergeOptions(options);
        return stackIntercept<ShareProfileRequest, ShareProfileResponse>("unary", this._transport, method, opt, input);
    }
}
