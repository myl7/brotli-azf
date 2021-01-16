import logging

import azure.functions as func
import brotli


def main(req: func.HttpRequest) -> func.HttpResponse:
    src = req.get_body()
    logging.info(f'get {len(src)} bytes')

    try:
        res = brotli.compress(src)
    except brotli.error:
        logging.error(f'enc {len(src)} bytes failed')
        return func.HttpResponse('brotli enc failed', status_code=400)

    logging.info(f'enc {len(src)} bytes ok')
    return func.HttpResponse(res, mimetype='application/octet-stream')
