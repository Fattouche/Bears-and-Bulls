from django.shortcuts import render

from django.http import HttpResponse

import django_rq

from .models import *


def index(request):
    return HttpResponse("Init view of stock exchange index")


def add(request):
    return HttpResponse("Init view of stock exchange index")


def quote(request):
    # placeholders
    symbol = "abc"
    user_id = "1234"
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
