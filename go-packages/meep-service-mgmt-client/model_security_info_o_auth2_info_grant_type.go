/*
 * Copyright (c) 2020  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the \"License\");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an \"AS IS\" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * AdvantEDGE MEC Service Management API
 *
 * The ETSI MEC ISG MEC011 MEC Service Management API described using OpenAPI
 *
 * API version: 2.1.1
 * Contact: cti_support@etsi.org
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

// SecurityInfoOAuth2InfoGrantType : OAuth 2.0 grant type
type SecurityInfoOAuth2InfoGrantType string

// List of SecurityInfo.OAuth2Info.GrantType
const (
	AUTHORIZATION_CODE_SecurityInfoOAuth2InfoGrantType SecurityInfoOAuth2InfoGrantType = "OAUTH2_AUTHORIZATION_CODE"
	IMPLICIT_GRANT_SecurityInfoOAuth2InfoGrantType     SecurityInfoOAuth2InfoGrantType = "OAUTH2_IMPLICIT_GRANT"
	RESOURCE_OWNER_SecurityInfoOAuth2InfoGrantType     SecurityInfoOAuth2InfoGrantType = "OAUTH2_RESOURCE_OWNER"
	CLIENT_CREDENTIALS_SecurityInfoOAuth2InfoGrantType SecurityInfoOAuth2InfoGrantType = "OAUTH2_CLIENT_CREDENTIALS"
)