var boiling = false;
var socket = null;

$(document).ready(function() {
    socket = io();

    socket.on('connect', function() {
        alert('connected!');
    });

    $('.fill-y').css({ height: $(window).height() });

    $('.btn.boil').click(function() {
        boil(this);
    });
});

function boil(btn) {
    boiling = true;
    $(btn).addClass('disabled');
    $(btn).children('i').html('grain').css({ animation: 'fadeIn 1s infinite alternate'});

    dots($(btn).children('span'), 'Boiling');
}

function dots(obj, str, n) {
    n = n || 0;
    obj.html(str + Array(n + 1).join('.') + Array(4 - n).join('&nbsp;'));
    if(boiling == true) {
        setTimeout(function() {
            dots(obj, str, (n < 3)? (n + 1): 0);
        }, 200);
    }
}