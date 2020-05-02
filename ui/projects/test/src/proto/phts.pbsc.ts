/* tslint:disable */
/* eslint-disable */
//
// THIS IS A GENERATED FILE
// DO NOT MODIFY IT! YOUR CHANGES WILL BE LOST
import { Inject, Injectable } from '@angular/core';
import {
  GrpcCallType,
  GrpcClient,
  GrpcClientFactory,
  GrpcClientSettings,
  GrpcEvent
} from '@ngx-grpc/common';
import {
  GRPC_CLIENT_FACTORY,
  GrpcHandler,
  takeMessages,
  throwStatusErrors
} from '@ngx-grpc/core';
import { Metadata } from 'grpc-web';
import { Observable } from 'rxjs';
import * as thisProto from './phts.pb';
import { GRPC_VERSION_SERVICE_CLIENT_SETTINGS } from './phts.pbconf';
@Injectable({
  providedIn: 'root'
})
export class VersionServiceClient {
  private client: GrpcClient;

  constructor(
    @Inject(GRPC_VERSION_SERVICE_CLIENT_SETTINGS) settings: GrpcClientSettings,
    @Inject(GRPC_CLIENT_FACTORY) clientFactory: GrpcClientFactory,
    private handler: GrpcHandler
  ) {
    this.client = clientFactory.createClient('VersionService', settings);
  }

  /**
   * Unary RPC. Emits messages and throws errors on non-zero status codes
   * @param thisProto.VersionRequest request
   * @param Metadata metadata
   * @return Observable<thisProto.VersionResponse>
   */
  get(
    requestData: thisProto.VersionRequest,
    requestMetadata: Metadata = {}
  ): Observable<thisProto.VersionResponse> {
    return this.get$eventStream(requestData, requestMetadata).pipe(
      throwStatusErrors(),
      takeMessages()
    );
  }

  /**
   * Unary RPC. Emits data and status events; does not throw errors by design
   * @param thisProto.VersionRequest request
   * @param Metadata metadata
   * @return Observable<GrpcEvent<thisProto.VersionResponse>>
   */
  get$eventStream(
    requestData: thisProto.VersionRequest,
    requestMetadata: Metadata = {}
  ): Observable<GrpcEvent<thisProto.VersionResponse>> {
    return this.handler.handle({
      type: GrpcCallType.unary,
      client: this.client,
      path: '/VersionService/Get',
      requestData,
      requestMetadata,
      requestClass: thisProto.VersionRequest,
      responseClass: thisProto.VersionResponse
    });
  }
}
