require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.js");
window.Bounce = require('bounce.js');
window.Noty = require('noty');
const ClassicEditor = require('@ckeditor/ckeditor5-build-classic');
$(() => {
    ClassicEditor.create(document.querySelector('#article-Content'));
    var count = Object.keys(flash).length;
    if (count > 0) {
        for (let k in flash) {
            flash[k].forEach(msg => {
                writeNoty(msg, k)
            });
        }
    }
});