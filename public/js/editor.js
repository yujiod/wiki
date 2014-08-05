var editor = null;
$(function () {
    // enable CodeMirror
    editor = CodeMirror.fromTextArea($('#Editor').get(0), {
        mode: "markdown",
        lineWrapping: true
    });
    $('.CodeMirror').attr('id', 'CodeMirror').addClass('tab-pane active');

    // live preview
    var preview = function () {
        var value = editor.getValue();
        while (value.match(/(\[\[([^\]\[\/]+)\]\])/)) {
            value = value.replace(RegExp.$1, '<a href="/page/'+encodeURIComponent(RegExp.$2)+'">'+RegExp.$2+'</a>');
        }
        $('#Preview').html(marked(value));
    };
    CodeMirror.on(editor, "change", preview);
    preview();

    // scroll sync
    $('.CodeMirror-scroll').on('scroll', function(e){
        var other = $('#Preview').get(0);
        var percentage = this.scrollTop / (this.scrollHeight - this.offsetHeight);
        other.scrollTop = percentage * (other.scrollHeight - other.offsetHeight);
    });

    // editor toolbox
    /*
    $('.btn-toolbar .btn[data-toggle!="dropdown"]').click(function () {
        $(this).blur();
        editor.focus();
    });
    $('.btn-bold').click(function () {
        var str = editor.doc.getSelection();
        editor.doc.replaceSelection('**'+str+'**');
    });
    $('.btn-italic').click(function () {
        var str = editor.doc.getSelection();
        editor.doc.replaceSelection('*'+str+'*');
    });
    var insertLineHead = function (stringHead, insertIndent) {
        var cursor = editor.doc.getCursor();
        var lineHandle = editor.doc.getLineHandle(cursor.line);
        if (lineHandle.text.replace(/(^\s+)|(\s+$)/g, "") === "") {
            editor.doc.setSelection(cursor);
            editor.doc.replaceSelection(stringHead);
        } else {
            editor.doc.setSelection({line: cursor.line, ch: 9999});
            var indent = '';
            if (insertIndent) {
                indent = lineHandle.text.replace(/^([\s]*).*$/, "$1");
            }
            editor.doc.replaceSelection("\n"+indent+stringHead);
        }
    };
    $('.btn-list-ul').click(function () {
        insertLineHead('* ', true);
    });
    $('.btn-list-ol').click(function () {
        insertLineHead('1. ', true);
    });
    $('.btn-header').click(function () {
        var headString = '';
        for (var i=1;i<=$(this).attr('data-header');i++) {
            headString += '#';
        }
        insertLineHead(headString+' ', false);
        $('.dropdown.open').removeClass('open');
        $('.btn-header').blur();
        editor.focus();
        return false;
    });
    $('.btn-chain').click(function () {
        var str = editor.doc.getSelection();
        editor.doc.replaceSelection('['+str+']()');
    });
    var replaceLineHead = function (stringHead) {
        var cursor = editor.doc.getCursor();
        var lineHandle = editor.doc.getLineHandle(cursor.line);
        if (lineHandle.text.replace(/(^\s+)|(\s+$)/g, "") === "") {
            editor.doc.setSelection(cursor);
            editor.doc.replaceSelection(stringHead);
        } else {
            editor.doc.setSelection({line: cursor.line, ch: 9999});
            var indent = '';
            if (insertIndent) {
                indent = lineHandle.text.replace(/^([\s]*).*$/, "$1");
            }
            editor.doc.replaceSelection("\n"+indent+stringHead);
        }
    };
    $('.btn-quote').click(function () {

    });
    $('.btn-code').click(function () {

    });
    */
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
});