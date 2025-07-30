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
import { DomainCodeFiles, DomainIndexResult, V1CliCreateParams } from "./types";

/**
 * @description 运行monkeycode-cli命令
 *
 * @tags CLI
 * @name V1CliCreate
 * @summary 运行monkeycode-cli命令
 * @request POST:/api/v1/cli/{command}
 * @response `200` `(DomainIndexResult)[]` 输出结果
 * @response `500` `WebResp` 内部错误
 */

export const v1CliCreate = (
  { command, ...query }: V1CliCreateParams,
  codeFiles: DomainCodeFiles,
  params: RequestParams = {},
) =>
  request<DomainIndexResult[]>({
    path: `/api/v1/cli/${command}`,
    method: "POST",
    query: query,
    body: codeFiles,
    type: ContentType.Json,
    format: "json",
    ...params,
  });
