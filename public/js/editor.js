var editor = null;
$(function () {
    // enable CodeMirror
    editor = CodeMirror.fromTextArea($('#Editor').get(0), {
        mode: "gfm",
        lineWrapping: true
    });
    $('.CodeMirror').attr('id', 'CodeMirror').addClass('tab-pane active');

    // live preview
    var preview = function () {
        var value = editor.getValue();
        while (value.match(/(\[\[([^\]\[\|]+)(\|([^\]\[]*))?\]\])/)) {
            var bracketLink = RegExp.$1;
            var title = RegExp.$2;
            var alias = RegExp.$4;
            if (alias === '') {
                alias = title;
            }
            value = value.replace(bracketLink, '<a href="/page/'+encodeURIComponent(alias.replace(/\//g, '-'))+'">'+title+'</a>');
        }
        $('#Preview').html(marked(value).replace(/<table>/g, '<table class="table table-bordered table-striped">'));
        $('body').trigger('modified');
    };
    CodeMirror.on(editor, "change", preview);
    preview();

    // scroll sync
    $('.CodeMirror-scroll').on('scroll', function(e){
        var other = $('#Preview').get(0);
        var percentage = this.scrollTop / (this.scrollHeight - this.offsetHeight);
        other.scrollTop = percentage * (other.scrollHeight - other.offsetHeight);
    });

    $('.btn-fullscreen').click(function () {
        var target = $('.editor').get(0);
        if (target.webkitRequestFullscreen) {
            // Chrome15+, Safari5.1+, Opera15+
            target.webkitRequestFullscreen();
        } else if (target.mozRequestFullScreen) {
            // Firefox 10+
            target.mozRequestFullScreen();
        } else if (target.msRequestFullscreen) {
            //IE11+
            target.msRequestFullscreen();
        } else if (target.requestFullscreen) {
            // HTML5 Fullscreen API
            target.requestFullscreen();
        }
        return false;
    });

    editor.focus();
});