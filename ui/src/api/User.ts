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
  DeleteDeleteAdminParams,
  DeleteDeleteUserParams,
  DomainAdminUser,
  DomainCreateAdminReq,
  DomainInviteResp,
  DomainListAdminLoginHistoryResp,
  DomainListAdminUserResp,
  DomainListLoginHistoryResp,
  DomainListUserResp,
  DomainLoginReq,
  DomainLoginResp,
  DomainRegisterReq,
  DomainSetting,
  DomainUpdateSettingReq,
  DomainUpdateUserReq,
  DomainUser,
  GetAdminLoginHistoryParams,
  GetListAdminUserParams,
  GetListUserParams,
  GetLoginHistoryParams,
  WebResp,
} from "./types";

/**
 * @description 创建管理员
 *
 * @tags User
 * @name PostCreateAdmin
 * @summary 创建管理员
 * @request POST:/api/v1/admin/create
 * @response `200` `(WebResp & {
    data?: DomainAdminUser,

})` OK
 */

export const postCreateAdmin = (
  param: DomainCreateAdminReq,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: DomainAdminUser;
    }
  >({
    path: `/api/v1/admin/create`,
    method: "POST",
    body: param,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 删除管理员
 *
 * @tags User
 * @name DeleteDeleteAdmin
 * @summary 删除管理员
 * @request DELETE:/api/v1/admin/delete
 * @response `200` `(WebResp & {
    data?: Record<string, any>,

})` OK
 */

export const deleteDeleteAdmin = (
  query: DeleteDeleteAdminParams,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: Record<string, any>;
    }
  >({
    path: `/api/v1/admin/delete`,
    method: "DELETE",
    query: query,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 获取管理员用户列表
 *
 * @tags User
 * @name GetListAdminUser
 * @summary 获取管理员用户列表
 * @request GET:/api/v1/admin/list
 * @response `200` `(WebResp & {
    data?: DomainListAdminUserResp,

})` OK
 */

export const getListAdminUser = (
  query: GetListAdminUserParams,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: DomainListAdminUserResp;
    }
  >({
    path: `/api/v1/admin/list`,
    method: "GET",
    query: query,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 管理员登录
 *
 * @tags User
 * @name PostAdminLogin
 * @summary 管理员登录
 * @request POST:/api/v1/admin/login
 * @response `200` `(WebResp & {
    data?: DomainAdminUser,

})` OK
 */

export const postAdminLogin = (
  param: DomainLoginReq,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: DomainAdminUser;
    }
  >({
    path: `/api/v1/admin/login`,
    method: "POST",
    body: param,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 获取管理员登录历史
 *
 * @tags User
 * @name GetAdminLoginHistory
 * @summary 获取管理员登录历史
 * @request GET:/api/v1/admin/login-history
 * @response `200` `(WebResp & {
    data?: DomainListAdminLoginHistoryResp,

})` OK
 */

export const getAdminLoginHistory = (
  query: GetAdminLoginHistoryParams,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: DomainListAdminLoginHistoryResp;
    }
  >({
    path: `/api/v1/admin/login-history`,
    method: "GET",
    query: query,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 获取系统设置
 *
 * @tags User
 * @name GetGetSetting
 * @summary 获取系统设置
 * @request GET:/api/v1/admin/setting
 * @response `200` `(WebResp & {
    data?: DomainSetting,

})` OK
 */

export const getGetSetting = (params: RequestParams = {}) =>
  request<
    WebResp & {
      data?: DomainSetting;
    }
  >({
    path: `/api/v1/admin/setting`,
    method: "GET",
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 更新系统设置
 *
 * @tags User
 * @name PutUpdateSetting
 * @summary 更新系统设置
 * @request PUT:/api/v1/admin/setting
 * @response `200` `(WebResp & {
    data?: DomainSetting,

})` OK
 */

export const putUpdateSetting = (
  param: DomainUpdateSettingReq,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: DomainSetting;
    }
  >({
    path: `/api/v1/admin/setting`,
    method: "PUT",
    body: param,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 下载VSCode插件
 *
 * @tags User
 * @name GetVsixDownload
 * @summary 下载VSCode插件
 * @request GET:/api/v1/static/vsix
 */

export const getVsixDownload = (params: RequestParams = {}) =>
  request<unknown>({
    path: `/api/v1/static/vsix`,
    method: "GET",
    type: ContentType.Json,
    ...params,
  });

/**
 * @description 删除用户
 *
 * @tags User
 * @name DeleteDeleteUser
 * @summary 删除用户
 * @request DELETE:/api/v1/user/delete
 * @response `200` `(WebResp & {
    data?: Record<string, any>,

})` OK
 */

export const deleteDeleteUser = (
  query: DeleteDeleteUserParams,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: Record<string, any>;
    }
  >({
    path: `/api/v1/user/delete`,
    method: "DELETE",
    query: query,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 获取用户邀请码
 *
 * @tags User
 * @name GetInvite
 * @summary 获取用户邀请码
 * @request GET:/api/v1/user/invite
 * @response `200` `(WebResp & {
    data?: DomainInviteResp,

})` OK
 */

export const getInvite = (params: RequestParams = {}) =>
  request<
    WebResp & {
      data?: DomainInviteResp;
    }
  >({
    path: `/api/v1/user/invite`,
    method: "GET",
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 获取用户列表
 *
 * @tags User
 * @name GetListUser
 * @summary 获取用户列表
 * @request GET:/api/v1/user/list
 * @response `200` `(WebResp & {
    data?: DomainListUserResp,

})` OK
 */

export const getListUser = (
  query: GetListUserParams,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: DomainListUserResp;
    }
  >({
    path: `/api/v1/user/list`,
    method: "GET",
    query: query,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 用户登录
 *
 * @tags User
 * @name PostLogin
 * @summary 用户登录
 * @request POST:/api/v1/user/login
 * @response `200` `(WebResp & {
    data?: DomainLoginResp,

})` OK
 */

export const postLogin = (param: DomainLoginReq, params: RequestParams = {}) =>
  request<
    WebResp & {
      data?: DomainLoginResp;
    }
  >({
    path: `/api/v1/user/login`,
    method: "POST",
    body: param,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 获取用户登录历史
 *
 * @tags User
 * @name GetLoginHistory
 * @summary 获取用户登录历史
 * @request GET:/api/v1/user/login-history
 * @response `200` `(WebResp & {
    data?: DomainListLoginHistoryResp,

})` OK
 */

export const getLoginHistory = (
  query: GetLoginHistoryParams,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: DomainListLoginHistoryResp;
    }
  >({
    path: `/api/v1/user/login-history`,
    method: "GET",
    query: query,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 注册用户
 *
 * @tags User
 * @name PostRegister
 * @summary 注册用户
 * @request POST:/api/v1/user/register
 * @response `200` `(WebResp & {
    data?: DomainUser,

})` OK
 */

export const postRegister = (
  param: DomainRegisterReq,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: DomainUser;
    }
  >({
    path: `/api/v1/user/register`,
    method: "POST",
    body: param,
    type: ContentType.Json,
    format: "json",
    ...params,
  });

/**
 * @description 更新用户
 *
 * @tags User
 * @name PutUpdateUser
 * @summary 更新用户
 * @request PUT:/api/v1/user/update
 * @response `200` `(WebResp & {
    data?: DomainUser,

})` OK
 */

export const putUpdateUser = (
  param: DomainUpdateUserReq,
  params: RequestParams = {},
) =>
  request<
    WebResp & {
      data?: DomainUser;
    }
  >({
    path: `/api/v1/user/update`,
    method: "PUT",
    body: param,
    type: ContentType.Json,
    format: "json",
    ...params,
  });
