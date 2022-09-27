// @generated by protobuf-ts 2.8.1 with parameter long_type_string,generate_dependencies
// @generated from protobuf file "share/share.proto" (package "polarsignals.share", syntax proto3)
// tslint:disable
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { ShareService } from "./share";
import type { ProfileTypesResponse } from "./share";
import type { ProfileTypesRequest } from "./share";
import type { QueryResponse } from "./share";
import type { QueryRequest } from "./share";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { UploadResponse } from "./share";
import type { UploadRequest } from "./share";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * Service that exposes APIs for sharing profiles.
 *
 * @generated from protobuf service parca.share.v1alpha1.ShareService
 */
export interface IShareServiceClient {
    /**
     * Uploads the profile and returns the link that can be used to access it.
     *
     * @generated from protobuf rpc: Upload(parca.share.v1alpha1.UploadRequest) returns (parca.share.v1alpha1.UploadResponse);
     */
    upload(input: UploadRequest, options?: RpcOptions): UnaryCall<UploadRequest, UploadResponse>;
    /**
     * Query performs a profile query
     *
     * @generated from protobuf rpc: Query(parca.share.v1alpha1.QueryRequest) returns (parca.share.v1alpha1.QueryResponse);
     */
    query(input: QueryRequest, options?: RpcOptions): UnaryCall<QueryRequest, QueryResponse>;
    /**
     * ProfileTypes returns the list of available profile types.
     *
     * @generated from protobuf rpc: ProfileTypes(parca.share.v1alpha1.ProfileTypesRequest) returns (parca.share.v1alpha1.ProfileTypesResponse);
     */
    profileTypes(input: ProfileTypesRequest, options?: RpcOptions): UnaryCall<ProfileTypesRequest, ProfileTypesResponse>;
}
/**
 * Service that exposes APIs for sharing profiles.
 *
 * @generated from protobuf service parca.share.v1alpha1.ShareService
 */
export class ShareServiceClient implements IShareServiceClient, ServiceInfo {
    typeName = ShareService.typeName;
    methods = ShareService.methods;
    options = ShareService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * Uploads the profile and returns the link that can be used to access it.
     *
     * @generated from protobuf rpc: Upload(parca.share.v1alpha1.UploadRequest) returns (parca.share.v1alpha1.UploadResponse);
     */
    upload(input: UploadRequest, options?: RpcOptions): UnaryCall<UploadRequest, UploadResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<UploadRequest, UploadResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * Query performs a profile query
     *
     * @generated from protobuf rpc: Query(parca.share.v1alpha1.QueryRequest) returns (parca.share.v1alpha1.QueryResponse);
     */
    query(input: QueryRequest, options?: RpcOptions): UnaryCall<QueryRequest, QueryResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<QueryRequest, QueryResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * ProfileTypes returns the list of available profile types.
     *
     * @generated from protobuf rpc: ProfileTypes(parca.share.v1alpha1.ProfileTypesRequest) returns (parca.share.v1alpha1.ProfileTypesResponse);
     */
    profileTypes(input: ProfileTypesRequest, options?: RpcOptions): UnaryCall<ProfileTypesRequest, ProfileTypesResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<ProfileTypesRequest, ProfileTypesResponse>("unary", this._transport, method, opt, input);
    }
}
