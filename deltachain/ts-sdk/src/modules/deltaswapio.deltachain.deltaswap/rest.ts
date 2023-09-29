//@ts-nocheck
/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface ProtobufAny {
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

/**
* message SomeRequest {
         Foo some_parameter = 1;
         PageRequest pagination = 2;
 }
*/
export interface V1Beta1PageRequest {
  /**
   * key is a value returned in PageResponse.next_key to begin
   * querying the next page most efficiently. Only one of offset or key
   * should be set.
   * @format byte
   */
  key?: string;

  /**
   * offset is a numeric offset that can be used when key is unavailable.
   * It is less efficient than using key. Only one of offset or key should
   * be set.
   * @format uint64
   */
  offset?: string;

  /**
   * limit is the total number of results to be returned in the result page.
   * If left empty it will default to a value to be set by each app.
   * @format uint64
   */
  limit?: string;

  /**
   * count_total is set to true  to indicate that the result set should include
   * a count of the total number of items available for pagination in UIs.
   * count_total is only respected when offset is used. It is ignored when key
   * is set.
   */
  count_total?: boolean;

  /**
   * reverse is set to true if results are to be returned in the descending order.
   *
   * Since: cosmos-sdk 0.43
   */
  reverse?: boolean;
}

/**
* PageResponse is to be embedded in gRPC response messages where the
corresponding request message has used PageRequest.

 message SomeResponse {
         repeated Bar results = 1;
         PageResponse page = 2;
 }
*/
export interface V1Beta1PageResponse {
  /** @format byte */
  next_key?: string;

  /** @format uint64 */
  total?: string;
}

export interface DeltaswapConfig {
  /** @format uint64 */
  phylax_set_expiration?: string;

  /** @format byte */
  governance_emitter?: string;

  /** @format int64 */
  governance_chain?: number;

  /** @format int64 */
  chain_id?: number;
}

export interface DeltaswapConsensusPhylaxSetIndex {
  /** @format int64 */
  index?: number;
}

export interface DeltaswapPhylaxSet {
  /** @format int64 */
  index?: number;
  keys?: string[];

  /** @format uint64 */
  expirationTime?: string;
}

export interface DeltaswapPhylaxValidator {
  /** @format byte */
  phylaxKey?: string;

  /** @format byte */
  validatorAddr?: string;
}

export type DeltaswapMsgExecuteGovernanceVAAResponse = object;

export interface DeltaswapMsgInstantiateContractResponse {
  /** Address is the bech32 address of the new contract instance. */
  address?: string;

  /** @format byte */
  data?: string;
}

export type DeltaswapMsgRegisterAccountAsPhylaxResponse = object;

export interface DeltaswapMsgStoreCodeResponse {
  /** @format uint64 */
  code_id?: string;
}

export interface DeltaswapQueryAllPhylaxSetResponse {
  PhylaxSet?: DeltaswapPhylaxSet[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

export interface DeltaswapQueryAllPhylaxValidatorResponse {
  phylaxValidator?: DeltaswapPhylaxValidator[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

export interface DeltaswapQueryAllReplayProtectionResponse {
  replayProtection?: DeltaswapReplayProtection[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

export interface DeltaswapQueryAllSequenceCounterResponse {
  sequenceCounter?: DeltaswapSequenceCounter[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

export interface DeltaswapQueryGetConfigResponse {
  Config?: DeltaswapConfig;
}

export interface DeltaswapQueryGetConsensusPhylaxSetIndexResponse {
  ConsensusPhylaxSetIndex?: DeltaswapConsensusPhylaxSetIndex;
}

export interface DeltaswapQueryGetPhylaxSetResponse {
  PhylaxSet?: DeltaswapPhylaxSet;
}

export interface DeltaswapQueryGetPhylaxValidatorResponse {
  phylaxValidator?: DeltaswapPhylaxValidator;
}

export interface DeltaswapQueryGetReplayProtectionResponse {
  replayProtection?: DeltaswapReplayProtection;
}

export interface DeltaswapQueryGetSequenceCounterResponse {
  sequenceCounter?: DeltaswapSequenceCounter;
}

export interface DeltaswapQueryLatestPhylaxSetIndexResponse {
  /** @format int64 */
  latestPhylaxSetIndex?: number;
}

export interface DeltaswapReplayProtection {
  index?: string;
}

export interface DeltaswapSequenceCounter {
  index?: string;

  /** @format uint64 */
  sequence?: string;
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: keyof Omit<Body, "body" | "bodyUsed">;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (securityData: SecurityDataType) => RequestParams | void;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "";
  private securityData: SecurityDataType = null as any;
  private securityWorker: null | ApiConfig<SecurityDataType>["securityWorker"] = null;
  private abortControllers = new Map<CancelToken, AbortController>();

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType) => {
    this.securityData = data;
  };

  private addQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];

    return (
      encodeURIComponent(key) +
      "=" +
      encodeURIComponent(Array.isArray(value) ? value.join(",") : typeof value === "number" ? value : `${value}`)
    );
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter((key) => "undefined" !== typeof query[key]);
    return keys
      .map((key) =>
        typeof query[key] === "object" && !Array.isArray(query[key])
          ? this.toQueryString(query[key] as QueryParamsType)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string") ? JSON.stringify(input) : input,
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((data, key) => {
        data.append(key, input[key]);
        return data;
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  private mergeRequestParams(params1: RequestParams, params2?: RequestParams): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  private createAbortSignal = (cancelToken: CancelToken): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format = "json",
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams = (secure && this.securityWorker && this.securityWorker(this.securityData)) || {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];

    return fetch(`${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`, {
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      signal: cancelToken ? this.createAbortSignal(cancelToken) : void 0,
      body: typeof body === "undefined" || body === null ? null : payloadFormatter(body),
    }).then(async (response) => {
      const r = response as HttpResponse<T, E>;
      r.data = (null as unknown) as T;
      r.error = (null as unknown) as E;

      const data = await response[format]()
        .then((data) => {
          if (r.ok) {
            r.data = data;
          } else {
            r.error = data;
          }
          return r;
        })
        .catch((e) => {
          r.error = e;
          return r;
        });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title deltaswap/config.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryConfig
   * @summary Queries a config by index.
   * @request GET:/deltaswapio/deltachain/deltaswap/config
   */
  queryConfig = (params: RequestParams = {}) =>
    this.request<DeltaswapQueryGetConfigResponse, RpcStatus>({
      path: `/deltaswapio/deltachain/deltaswap/config`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryConsensusPhylaxSetIndex
   * @summary Queries a ConsensusPhylaxSetIndex by index.
   * @request GET:/deltaswapio/deltachain/deltaswap/consensus_phylax_set_index
   */
  queryConsensusPhylaxSetIndex = (params: RequestParams = {}) =>
    this.request<DeltaswapQueryGetConsensusPhylaxSetIndexResponse, RpcStatus>({
      path: `/deltaswapio/deltachain/deltaswap/consensus_phylax_set_index`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryPhylaxSetAll
   * @summary Queries a list of phylaxSet items.
   * @request GET:/deltaswapio/deltachain/deltaswap/phylaxSet
   */
  queryPhylaxSetAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<DeltaswapQueryAllPhylaxSetResponse, RpcStatus>({
      path: `/deltaswapio/deltachain/deltaswap/phylaxSet`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryPhylaxSet
   * @summary Queries a phylaxSet by index.
   * @request GET:/deltaswapio/deltachain/deltaswap/phylaxSet/{index}
   */
  queryPhylaxSet = (index: number, params: RequestParams = {}) =>
    this.request<DeltaswapQueryGetPhylaxSetResponse, RpcStatus>({
      path: `/deltaswapio/deltachain/deltaswap/phylaxSet/${index}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryPhylaxValidatorAll
   * @summary Queries a list of PhylaxValidator items.
   * @request GET:/deltaswapio/deltachain/deltaswap/phylax_validator
   */
  queryPhylaxValidatorAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<DeltaswapQueryAllPhylaxValidatorResponse, RpcStatus>({
      path: `/deltaswapio/deltachain/deltaswap/phylax_validator`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryPhylaxValidator
   * @summary Queries a PhylaxValidator by index.
   * @request GET:/deltaswapio/deltachain/deltaswap/phylax_validator/{phylaxKey}
   */
  queryPhylaxValidator = (phylaxKey: string, params: RequestParams = {}) =>
    this.request<DeltaswapQueryGetPhylaxValidatorResponse, RpcStatus>({
      path: `/deltaswapio/deltachain/deltaswap/phylax_validator/${phylaxKey}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryLatestPhylaxSetIndex
   * @summary Queries a list of LatestPhylaxSetIndex items.
   * @request GET:/deltaswapio/deltachain/deltaswap/latest_phylax_set_index
   */
  queryLatestPhylaxSetIndex = (params: RequestParams = {}) =>
    this.request<DeltaswapQueryLatestPhylaxSetIndexResponse, RpcStatus>({
      path: `/deltaswapio/deltachain/deltaswap/latest_phylax_set_index`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryReplayProtectionAll
   * @summary Queries a list of replayProtection items.
   * @request GET:/deltaswapio/deltachain/deltaswap/replayProtection
   */
  queryReplayProtectionAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<DeltaswapQueryAllReplayProtectionResponse, RpcStatus>({
      path: `/deltaswapio/deltachain/deltaswap/replayProtection`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryReplayProtection
   * @summary Queries a replayProtection by index.
   * @request GET:/deltaswapio/deltachain/deltaswap/replayProtection/{index}
   */
  queryReplayProtection = (index: string, params: RequestParams = {}) =>
    this.request<DeltaswapQueryGetReplayProtectionResponse, RpcStatus>({
      path: `/deltaswapio/deltachain/deltaswap/replayProtection/${index}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QuerySequenceCounterAll
   * @summary Queries a list of sequenceCounter items.
   * @request GET:/deltaswapio/deltachain/deltaswap/sequenceCounter
   */
  querySequenceCounterAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<DeltaswapQueryAllSequenceCounterResponse, RpcStatus>({
      path: `/deltaswapio/deltachain/deltaswap/sequenceCounter`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QuerySequenceCounter
   * @summary Queries a sequenceCounter by index.
   * @request GET:/deltaswapio/deltachain/deltaswap/sequenceCounter/{index}
   */
  querySequenceCounter = (index: string, params: RequestParams = {}) =>
    this.request<DeltaswapQueryGetSequenceCounterResponse, RpcStatus>({
      path: `/deltaswapio/deltachain/deltaswap/sequenceCounter/${index}`,
      method: "GET",
      format: "json",
      ...params,
    });
}
