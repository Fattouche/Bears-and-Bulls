from django.shortcuts import render

from django.http import HttpResponse

import django_rq

from .models import *

from decimal import Decimal


def index(request):
    return HttpResponse("Init view of stock exchange index")


def add(request):
    params = request.POST
    user_id = params.get('user_id')
    amount = params.get('amount')
    user = User.objects.get(user_id=user_id)
    user.balance += Decimal(amount)
    user.save()
    return HttpResponse("New balance {0}".format(user.balance), status=200)


def quote(request):
    # placeholders
    params = request.GET
    user_id = params.get('user_id')
    symbol = params.get('symbol')
    stock = cache.get(symbol)
    if(stock is None):
        price = execute_request(
            user_id=user_id, symbol=symbol, command="QUOTE")
        stock = Stock(symbol=symbol, price=price)
        cache.set(symbol, stock)
        django_rq.enqueue(stock.verify_triggers)

    return HttpResponse("stock: "+symbol)


def buy(request):
    return HttpResponse("Init view of stock exchange index")


def commit_buy(request):
    return HttpResponse("Init view of stock exchange index")


def cancel_buy(request):
    return HttpResponse("Init view of stock exchange index")


def sell(request):
    return HttpResponse("Init view of stock exchange index")


def commit_sell(request):
    return HttpResponse("Init view of stock exchange index")


def cancel_sell(request):
    return HttpResponse("Init view of stock exchange index")


def set_buy_amount(request):
    return HttpResponse("Init view of stock exchange index")


def set_sell_amount(request):
    return HttpResponse("Init view of stock exchange index")


def set_sell_trigger(request):
    return HttpResponse("Init view of stock exchange index")


def cancel_set_sell(request):
    return HttpResponse("Init view of stock exchange index")


def set_buy_trigger(request):
    return HttpResponse("Init view of stock exchange index")


def dumplog(request):
    return HttpResponse("Init view of stock exchange index")


def display_summary(request):
    return HttpResponse("Init view of stock exchange index")
