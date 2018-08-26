require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.js");
window.Bounce = require('bounce.js');
window.Noty = require('noty');
$(() => {
    var count = Object.keys(flash).length;
    if (count > 0) {
        for (let k in flash) {
            flash[k].forEach(msg => {
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
            });
        }
    }
});
