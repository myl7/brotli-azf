import logging

import azure.functions as func
import brotli


def main(req: func.HttpRequest) -> func.HttpResponse:
    src = req.get_body()
    logging.info(f'get {len(src)} bytes')

    try:
        res = brotli.decompress(src)
    except brotli.error:
        logging.error(f'dec {len(src)} bytes failed')
        return func.HttpResponse('brotli dec failed', status_code=400)

    logging.info(f'dec {len(src)} bytes ok')
    return func.HttpResponse(res, mimetype='application/octet-stream')
