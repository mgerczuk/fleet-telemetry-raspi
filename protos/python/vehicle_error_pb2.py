# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: vehicle_error.proto
# Protobuf Python Version: 5.28.3
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    28,
    3,
    '',
    'vehicle_error.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x13vehicle_error.proto\x12\x17telemetry.vehicle_error\x1a\x1fgoogle/protobuf/timestamp.proto\"\x83\x01\n\rVehicleErrors\x12\x35\n\x06\x65rrors\x18\x01 \x03(\x0b\x32%.telemetry.vehicle_error.VehicleError\x12.\n\ncreated_at\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x0b\n\x03vin\x18\x03 \x01(\t\"\xc6\x01\n\x0cVehicleError\x12.\n\ncreated_at\x18\x01 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x0c\n\x04name\x18\x02 \x01(\t\x12=\n\x04tags\x18\x03 \x03(\x0b\x32/.telemetry.vehicle_error.VehicleError.TagsEntry\x12\x0c\n\x04\x62ody\x18\x04 \x01(\t\x1a+\n\tTagsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\x42/Z-github.com/teslamotors/fleet-telemetry/protosb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'vehicle_error_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z-github.com/teslamotors/fleet-telemetry/protos'
  _globals['_VEHICLEERROR_TAGSENTRY']._loaded_options = None
  _globals['_VEHICLEERROR_TAGSENTRY']._serialized_options = b'8\001'
  _globals['_VEHICLEERRORS']._serialized_start=82
  _globals['_VEHICLEERRORS']._serialized_end=213
  _globals['_VEHICLEERROR']._serialized_start=216
  _globals['_VEHICLEERROR']._serialized_end=414
  _globals['_VEHICLEERROR_TAGSENTRY']._serialized_start=371
  _globals['_VEHICLEERROR_TAGSENTRY']._serialized_end=414
# @@protoc_insertion_point(module_scope)
