"""
rpc.py
------

This module implements the services required by the master `.proto`.
"""

import logging
from logging.config import fileConfig
import time
from concurrent import futures

import grpc

from common.utils import get_logger_cfg_fpath
from extractor import IMDb_Queries_Manager
from extractor.config import ExtractorConfig
import imdb_pb2
import imdb_pb2_grpc

try:
    fileConfig(get_logger_cfg_fpath())
except FileNotFoundError as e:
    print(e)
logger = logging.getLogger(__name__)


class ExtractionService(imdb_pb2_grpc.ExtractorServiceServicer):
    def __init__(self, *args, **kwargs):
        self.mgr = IMDb_Queries_Manager(ExtractorConfig())

    def InitiateExtraction(self, request, context):
        logger.info("Requesting to extract `{}`".format(request.item_name))
        self.mgr.add_query(request.item_name)
        successful = self.mgr.api_execute()
        logger.info("Request `{}` finished; success status: {}"
                    .format(request.item_name, successful))
        return imdb_pb2.ExtractionResponse(
            item_name=request.item_name, successful=successful)


def serve():
    PORT = ":8989"
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    imdb_pb2_grpc.add_ExtractorServiceServicer_to_server(
        ExtractionService(), server)
    server.add_insecure_port('[::]{}'.format(PORT))
    server.start()
    logger.info("Server started at port{}".format(PORT))
    try:
        while True:
            time.sleep(60 * 60 * 12)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()