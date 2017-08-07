var boiling = false;
var socket = null;
var room = 'jug';

$(document).ready(function() {
    socket = io();

    socket.on('connect', function() {
        console.log('connected.');
    })

    socket.on('boiled', onBoiled);
    socket.on('boiling', disableButton);

    socket.on('message', function(msg) {
        console.log(msg);
    })

    $('.fill-y').css({ height: $(window).height() });

    $('.btn.boil').click(boil);
});

function onBoiled() {
    console.log('boiled.');
    boiling = false;
    $('.btn.boil')
        .removeClass('disabled')
        .children('i')
        .html('local_drink')
        .css({ animation: 'none'});
    
    setTimeout(function() {
        location.reload();
    }, 1000);
}

function boil() {
    console.log('boiling..');
    console.log(socket.emit('boil'));
    disableButton();
}

function disableButton() {
    boiling = true;
    $('.btn.boil')
        .addClass('disabled')
        .children('i').html('grain')
        .css({ animation: 'fadeIn 1s infinite alternate'});

    dots($('.btn.boil').children('span'), 'Boiling');
}

function dots(obj, str, n) {
    n = n || 0;
    obj.html(str + Array(n + 1).join('.') + Array(4 - n).join('&nbsp;'));
    if(boiling == true) {
        setTimeout(function() {
            dots(obj, str, (n < 3)? (n + 1): 0);
        }, 200);
    } else {
        $('.btn.boil')
            .children('span')
            .html('Boil!');
    }
}