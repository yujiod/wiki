$(function () {
    $('button.click-to-disabled, input.click-to-disabled').click(function () {
        $(this).addClass('disabled').attr('disabled', 'disabled');
        $(this).parents('.btn-group').find('.btn').addClass('disabled').attr('disabled', 'disabled');
        $(this).parents('form').submit();
    });

    $('body').on('modified', function () {
        $('.contents div.highlight').each(function () {
            var lang = '';
            if ($(this).attr('class').match(/highlight-([\w]+)/)) {
                lang = RegExp.$1;
            }
            $(this).find('pre').wrapInner('<code></code>').find('code').addClass(lang);
        });
        $('.contents pre code[class*=lang-]').each(function () {
            if ($(this).attr('class').match(/lang-([\w]+)/)) {
                var lang = RegExp.$1;
                $(this).addClass(lang);
            }
        });
        $('pre code').each(function (i, block) {
            hljs.highlightBlock(block);
        });
    });
    $('body').trigger('modified');
});