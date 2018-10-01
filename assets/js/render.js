const moment = require('moment')
function writeData(data) {
    $("#result").empty();
    let pom = '';
    data['articles'].forEach(val => {
        let content = val.content.replace(/<\/?[^>]+(>|$)/g, "");
        let dots = "";
        if (content.length > 451) {
            dots = "..."
        }
        console.log(val.created_at)
        pom += '' +
            '        <a class="a_hover"  href="/article-detail/'+ val.id +'">'+
            '            <div class="row article">'+
            '                <div class="col-md-4">'+
            '                    <img src="/assets/' + val.title_image +'" height="200px" alt="">'+
            '                </div>'+
            '                <div class="col-md-8">'+
            '                    <div class="text">'+
            '                        <h1>'+ val.title +'</h1>'+
            '                        <p>'+ content.slice(0, 451) + dots +'</p>'+
            '                        <div class="time">'+
            '                           '+ moment(val.created_at).utcOffset('+0000').format("MMMM D, YYYY, HH:mm") +''+
            '                        </div>'+
            '                    </div>'+
            '                </div>'+
            '            </div>'+
            '        </a>';
    });
    $("#result").html(pom)
}

function writeNoty(msg, k) {
    new Noty({
        type: k,
        theme: 'mint',
        layout: 'bottomRight',
        timeout: 2000,
        text: msg,
        animation: {
            open: function (promise) {
                var n = this;
                new Bounce()
                    .translate({
                        from     : {x: 450, y: 0}, to: {x: 0, y: 0},
                        easing   : "bounce",
                        duration : 1000,
                        bounces  : 4,
                        stiffness: 3
                    })
                    .scale({
                        from     : {x: 1.2, y: 1}, to: {x: 1, y: 1},
                        easing   : "bounce",
                        duration : 1000,
                        delay    : 100,
                        bounces  : 4,
                        stiffness: 1
                    })
                    .scale({
                        from     : {x: 1, y: 1.2}, to: {x: 1, y: 1},
                        easing   : "bounce",
                        duration : 1000,
                        delay    : 100,
                        bounces  : 6,
                        stiffness: 1
                    })
                    .applyTo(n.barDom, {
                        onComplete: function () {
                            promise(function(resolve) {
                                resolve();
                            })
                        }
                    });
            },
            close: function (promise) {
                var n = this;
                new Bounce()
                    .translate({
                        from     : {x: 0, y: 0}, to: {x: 450, y: 0},
                        easing   : "bounce",
                        duration : 500,
                        bounces  : 4,
                        stiffness: 1
                    })
                    .applyTo(n.barDom, {
                        onComplete: function () {
                            promise(function(resolve) {
                                resolve();
                            })
                        }
                    });
            }
        }
    }).show();
}

window.writeNoty = writeNoty
window.writeData = writeData