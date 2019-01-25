from django.core.exceptions import ObjectDoesNotExist
from exchange.models import User

class CreateUserMiddleware(object):
    def __init__(self, get_response):
        self.get_response = get_response
    
    def __call__(self, request):
        query_dict = request.POST if request.method == 'POST' else request.GET
        
        if 'user_id' in query_dict:
            try:
                u_id = query_dict['user_id']
                user = User.objects.get(user_id=u_id)
            except ObjectDoesNotExist:
                User.objects.create(user_id=u_id, balance=0)
        
        return self.get_response(request)