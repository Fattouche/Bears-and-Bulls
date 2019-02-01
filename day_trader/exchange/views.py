from django.shortcuts import render

from django.http import HttpResponse

import django_rq

from .models.business_models import *
from django.core.cache import cache
from decimal import Decimal
from django.http import JsonResponse
import json
from .audit_logging import AuditLogger


def index(request):
    return HttpResponse("Init view of stock exchange index")


def add(request):
    params = json.loads(request.body)
    user_id = params.get('user_id')
    amount = Decimal(params.get('amount'))
    user = User.get(user_id)
    user.update_balance(amount)
    return JsonResponse({'action': 'add', 'balance': user.balance}, status=200)


def quote(request):
    params = request.GET
    user_id = params.get('user_id')
    symbol = params.get('symbol')
    stock = Stock.quote(symbol, user_id)
    return JsonResponse({'action': 'quote', 'stock': stock.symbol, 'price': stock.price}, status=200)


def buy(request):
    params = json.loads(request.body)
    user_id = params.get('user_id')
    symbol = params.get('symbol')
    amount = Decimal(params.get('amount'))
    user = User.get(user_id)
    err = user.perform_buy(symbol, amount)
    if err is not None:
        return JsonResponse({'action': 'buy', 'error': err}, status=400)
    return JsonResponse({'action': 'buy', 'balance': user.balance, 'stock': symbol, 'price': user.buy_stack[-1].purchase_price, 'amount': amount, 'valid_duration': '60'}, status=200)


def commit_buy(request):
    params = json.loads(request.body)
    user_id = params.get('user_id')
    user = User.get(user_id)
    buy = user.pop_from_buy_stack()
    if(buy is None):
        return JsonResponse({'action': 'commit_buy', 'error': 'no buy currently exists'}, status=404)
    if(is_expired(buy.timestamp)):
        return JsonResponse({'action': 'commit_buy', 'stock': buy.stock_symbol, 'error': 'buy has expired, please re-buy in order to commit'}, status=408)
    buy.commit(user)
    return JsonResponse({'action': 'commit_buy', 'stock': buy.stock_symbol, 'amount_bought': buy.stock_bought_amount, 'balance': user.balance}, status=200)


def cancel_buy(request):
    params = json.loads(request.body)
    user = User.get(params.get('user_id'))
    user.cancel_buy()
    return JsonResponse({'action': 'cancel_buy', 'balance': user.balance}, status=200)


def sell(request):
    params = json.loads(request.body)
    user_id = params.get('user_id')
    symbol = params.get('symbol')
    amount = Decimal(params.get('amount'))
    user = User.get(user_id)
    sell = user.pop_from_sell_stack()
    err = user.perform_sell(symbol, amount)
    if err is not None:
        return JsonResponse({'action': 'sell', 'error': err}, status=400)
    return JsonResponse({'action': 'sell', 'stock': symbol, 'price': user.sell_stack[-1].sell_price, 'amount': amount, 'valid_duration': '60'}, status=200)


def commit_sell(request):
    params = json.loads(request.body)
    user_id = params.get('user_id')
    user = User.get(user_id)
    sell = user.pop_from_sell_stack()
    if(sell is None):
        return JsonResponse({'action': 'commit_sell', 'error': 'no sell currently exists'}, status=404)
    if(is_expired(sell.timestamp)):
        return JsonResponse({'action': 'commit_sell', 'stock': sell.stock_symbol, 'error': 'sell has expired, please re-buy in order to commit'}, status=408)
    sell.commit(user)
    return JsonResponse({'action': 'commit_sell', 'stock': sell.stock_symbol, 'amount_sold': sell.stock_sold_amount, 'balance': user.balance}, status=200)


def cancel_sell(request):
    params = json.loads(request.body)
    user = User.get(params.get('user_id'))
    user.cancel_sell()
    return JsonResponse({'action': 'cancel_sell', 'balance': user.balance}, status=200)


def set_buy_amount(request):
    params = json.loads(request.body)
    user_id = params.get('user_id')
    symbol = params.get('symbol')
    amount = Decimal(params.get('amount'))
    user = User.get(user_id)
    if(user.balance < amount):
        return JsonResponse({'action': 'set_buy_amount', 'error': 'user balance too low'}, status=412)
    user.set_buy_amount(symbol, amount)
    return JsonResponse({'action': 'set_buy_amount'}, status=200)


def set_sell_amount(request):
    params = json.loads(request.body)
    user_id = params.get('user_id')
    symbol = params.get('symbol')
    amount = Decimal(params.get('amount'))
    user = User.get(user_id)
    user.set_sell_amount(symbol, amount)
    return JsonResponse({'action': 'set_sell_amount'}, status=200)


def set_buy_trigger(request):
    params = json.loads(request.body)
    user_id = params.get('user_id')
    symbol = params.get('symbol')
    amount = Decimal(params.get('amount'))
    user = User.get(user_id)
    trigger_set = user.set_buy_trigger(symbol, amount)
    if not trigger_set:
        return JsonResponse({'action': 'set_buy_trigger', 'error': 'no existing set buy amount'}, status=412)
    return JsonResponse({'action': 'set_buy_trigger'}, status=200)


def set_sell_trigger(request):
    params = json.loads(request.body)
    user_id = params.get('user_id')
    symbol = params.get('symbol')
    amount = Decimal(params.get('amount'))
    user = User.get(user_id)
    trigger_set = user.set_sell_trigger(symbol, amount)
    if not trigger_set:
        return JsonResponse({'action': 'set_sell_trigger', 'error': 'no existing set sell amount'}, status=412)
    return JsonResponse({'action': 'set_sell_trigger'}, status=200)


def cancel_set_buy(request):
    params = json.loads(request.body)
    user_id = params.get('user_id')
    symbol = params.get('symbol')
    user = User.get(user_id)
    set_buy_cancelled = user.cancel_set_buy(symbol)
    if not set_buy_cancelled:
        return JsonResponse({'action': 'cancel_set_buy', 'error': 'no set buy to cancel'}, status=412)
    return JsonResponse({'action': 'cancel_set_buy'}, status=200)


def cancel_set_sell(request):
    params = json.loads(request.body)
    user_id = params.get('user_id')
    symbol = params.get('symbol')
    user = User.get(user_id)
    set_sell_cancelled = user.cancel_set_sell(symbol)
    if not set_sell_cancelled:
        return JsonResponse({'action': 'cancel_set_sell', 'error': 'no set sell to cancel'}, status=412)
    return JsonResponse({'action': 'cancel_set_sell'}, status=200)


def dumplog(request):
    params = request.GET
    if 'filename' not in params:
        return JsonResponse({'error': 'filename is required'}, status=400)

    if 'user_id' not in params:
        AuditLogger.dump_system_logs(params['filename'])
    else:
        AuditLogger.dump_user_log(params['user_id'], params['filename'])

    return HttpResponse("System logs dumped to file {}".format(params['filename']))


def display_summary(request):
    return HttpResponse("Init view of stock exchange index")
