/* tslint:disable */
/* eslint-disable */
//
// THIS IS A GENERATED FILE
// DO NOT MODIFY IT! YOUR CHANGES WILL BE LOST
import { GrpcMessage, RecursivePartial } from '@ngx-grpc/common';
import { BinaryReader, BinaryWriter, ByteSource } from 'google-protobuf';
export class VersionResponse implements GrpcMessage {
  static toBinary(instance: VersionResponse) {
    const writer = new BinaryWriter();
    VersionResponse.toBinaryWriter(instance, writer);
    return writer.getResultBuffer();
  }

  static fromBinary(bytes: ByteSource) {
    const instance = new VersionResponse();
    VersionResponse.fromBinaryReader(instance, new BinaryReader(bytes));
    return instance;
  }

  static refineValues(instance: VersionResponse) {
    instance.sha = instance.sha || '';
    instance.buildTime = instance.buildTime || '';
  }

  static fromBinaryReader(instance: VersionResponse, reader: BinaryReader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) break;

      switch (reader.getFieldNumber()) {
        case 1:
          instance.sha = reader.readString();
          break;
        case 2:
          instance.buildTime = reader.readString();
          break;
        default:
          reader.skipField();
      }
    }

    VersionResponse.refineValues(instance);
  }

  static toBinaryWriter(instance: VersionResponse, writer: BinaryWriter) {
    if (instance.sha) {
      writer.writeString(1, instance.sha);
    }
    if (instance.buildTime) {
      writer.writeString(2, instance.buildTime);
    }
  }

  private _sha?: string;
  private _buildTime?: string;

  /**
   * Creates an object and applies default Protobuf values
   * @param VersionResponse value
   */
  constructor(value?: RecursivePartial<VersionResponse>) {
    value = value || {};
    this.sha = value.sha;
    this.buildTime = value.buildTime;
    VersionResponse.refineValues(this);
  }
  get sha(): string | undefined {
    return this._sha;
  }
  set sha(value: string | undefined) {
    this._sha = value;
  }
  get buildTime(): string | undefined {
    return this._buildTime;
  }
  set buildTime(value: string | undefined) {
    this._buildTime = value;
  }
  toObject() {
    return {
      sha: this.sha,
      buildTime: this.buildTime
    };
  }
  toJSON() {
    return this.toObject();
  }
}
export module VersionResponse {}
export class VersionRequest implements GrpcMessage {
  static toBinary(instance: VersionRequest) {
    const writer = new BinaryWriter();
    VersionRequest.toBinaryWriter(instance, writer);
    return writer.getResultBuffer();
  }

  static fromBinary(bytes: ByteSource) {
    const instance = new VersionRequest();
    VersionRequest.fromBinaryReader(instance, new BinaryReader(bytes));
    return instance;
  }

  static refineValues(instance: VersionRequest) {}

  static fromBinaryReader(instance: VersionRequest, reader: BinaryReader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) break;

      switch (reader.getFieldNumber()) {
        default:
          reader.skipField();
      }
    }

    VersionRequest.refineValues(instance);
  }

  static toBinaryWriter(instance: VersionRequest, writer: BinaryWriter) {}

  /**
   * Creates an object and applies default Protobuf values
   * @param VersionRequest value
   */
  constructor(value?: RecursivePartial<VersionRequest>) {
    value = value || {};
    VersionRequest.refineValues(this);
  }
  toObject() {
    return {};
  }
  toJSON() {
    return this.toObject();
  }
}
export module VersionRequest {}
