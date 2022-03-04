# minio

#### minio-version: [RELEASE.2022-02-26T02-54-46Z](https://github.com/minio/minio/tree/bc33db9fc0d30428b8240b88f7acaa5b03743ba2)



- ### 纠删码

  MinIO 使用纠删码来校验和保护数据免受硬件故障和静默数据损坏（Bit Rot）。使用最高级别的冗余，即使丢失多达一半 (N/2) 的驱动器，但仍然能够恢复数据。

  默认情况下，MinIO 在 N/2 个数据和 N/2 个奇偶校验驱动器上对对象进行分片。即数据盘 ( DataDrives )和 冗余盘 ( ParityDrives ) 个数相同， 所以我们真正可用的存储空间，只有我们总空间的一半大小。

  

  在下面的16个驱动器示例中，丢失任意的8个驱动器，仍可以从剩下的盘中进行数据恢复.。

  ![ErasureCode](ErasureCode.png)

- ### 文档

  - [官方文档(英文)](https://docs.min.io/)
  - [中文文档(旧版)](http://docs.minio.org.cn/docs/)
  - [Minio Cluster 详解](https://jicki.cn/minio-cluster/#minio-%E4%BB%8B%E7%BB%8D)