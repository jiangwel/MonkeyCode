/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

import request, { ContentType, RequestParams } from "./httpClient";
import {
  DomainAcceptCompletionReq,
  DomainModelListResp,
  WebResp,
} from "./types";

/**
 * @description 处理聊天补全请求
 *
 * @tags OpenAIV1
 * @name PostChatCompletion
 * @summary 处理聊天补全请求
 * @request POST:/v1/chat/completions
 * @response `200` `WebResp` OK
 */

export const postChatCompletion = (params: RequestParams = {}) =>
  request<WebResp>({
    path: `/v1/chat/completions`,
    method: "POST",
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 接受补全请求
 *
 * @tags OpenAIV1
 * @name PostAcceptCompletion
 * @summary 接受补全请求
 * @request POST:/v1/completion/accept
 * @response `200` `WebResp` OK
 */

export const postAcceptCompletion = (
  param: DomainAcceptCompletionReq,
  params: RequestParams = {},
) =>
  request<WebResp>({
    path: `/v1/completion/accept`,
    method: "POST",
    body: param,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 处理文本补全请求
 *
 * @tags OpenAIV1
 * @name PostCompletions
 * @summary 处理文本补全请求
 * @request POST:/v1/completions
 * @response `200` `WebResp` OK
 */

export const postCompletions = (params: RequestParams = {}) =>
  request<WebResp>({
    path: `/v1/completions`,
    method: "POST",
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 处理嵌入请求
 *
 * @tags OpenAIV1
 * @name PostEmbeddings
 * @summary 处理嵌入请求
 * @request POST:/v1/embeddings
 * @response `200` `WebResp` OK
 */

export const postEmbeddings = (params: RequestParams = {}) =>
  request<WebResp>({
    path: `/v1/embeddings`,
    method: "POST",
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 模型列表
 *
 * @tags OpenAIV1
 * @name GetModelList
 * @summary 模型列表
 * @request GET:/v1/models
 * @response `200` `(WebResp & {
    data?: DomainModelListResp,

})` OK
 */

export const getModelList = (params: RequestParams = {}) =>
  request<
    WebResp & {
      data?: DomainModelListResp;
    }
  >({
    path: `/v1/models`,
    method: "GET",
    type: ContentType.Json,
    format: "json",
    ...params,
  });
