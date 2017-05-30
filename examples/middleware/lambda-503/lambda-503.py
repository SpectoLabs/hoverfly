from time import sleep
from random import randint

def lambda_handler(event, context):
    sleep(randint(0,2))
 
    event['response'] = {}
    event['response']['status'] = 503
    event['response']['body'] = "PGh0bWw+DQo8aGVhZD48dGl0bGU+NTAwIEludGVybmFsIFNlcnZlciBFcnJvcjwvdGl0bGU+PC9oZWFkPg0KPGJvZHkgYmdjb2xvcj0id2hpdGUiPg0KPGNlbnRlcj48aDE+NTAwIEludGVybmFsIFNlcnZlciBFcnJvcjwvaDE+PC9jZW50ZXI+DQo8aHI+PGNlbnRlcj5uZ2lueC8xLjYuMjwvY2VudGVyPg0KPC9ib2R5Pg0KPC9odG1sPg=="
    event['response']['encodedBody'] = True
    event['response']['headers'] = {}

    
    return event
