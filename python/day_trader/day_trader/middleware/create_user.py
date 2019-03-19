from django.core.exceptions import ObjectDoesNotExist
from exchange.models import User
import json

class CreateUserMiddleware(object):
    def __init__(self, get_response):
        self.get_response = get_response
    
    def __call__(self, request):
        query_dict = request.GET if request.method == 'GET' else json.loads(request.body)

        if 'user_id' in query_dict:
            User.objects.get_or_create(
                user_id=query_dict['user_id'],
                defaults={'balance': 0},
            )

        return self.get_response(request)