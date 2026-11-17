"""Object storage for source documents and results (S3, or MinIO locally).

Keys are laid out as <tenant>/<document id>/<name>, so everything that
belongs to one document can be deleted by prefix.
"""

import asyncio

import boto3

SIGNED_URL_SECONDS = 900


class BlobStore:
    def __init__(self, bucket: str, endpoint: str | None):
        self.bucket = bucket
        self.s3 = boto3.client("s3", endpoint_url=endpoint)

    async def put(self, key: str, data: bytes, content_type: str) -> None:
        await asyncio.to_thread(
            self.s3.put_object, Bucket=self.bucket, Key=key, Body=data, ContentType=content_type
        )

    async def get(self, key: str) -> bytes:
        obj = await asyncio.to_thread(self.s3.get_object, Bucket=self.bucket, Key=key)
        return obj["Body"].read()

    async def delete_prefix(self, prefix: str) -> int:
        listed = await asyncio.to_thread(
            self.s3.list_objects_v2, Bucket=self.bucket, Prefix=prefix
        )
        keys = [{"Key": o["Key"]} for o in listed.get("Contents", [])]
        if keys:
            await asyncio.to_thread(
                self.s3.delete_objects, Bucket=self.bucket, Delete={"Objects": keys}
            )
        return len(keys)

    def signed_url(self, key: str) -> str:
        return self.s3.generate_presigned_url(
            "get_object",
            Params={"Bucket": self.bucket, "Key": key},
            ExpiresIn=SIGNED_URL_SECONDS,
        )
