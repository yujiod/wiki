$(function () {
    $('button.click-to-disabled, input.click-to-disabled').click(function () {
        $(this).addClass('disabled').attr('disabled', 'disabled');
        $(this).parents('.btn-group').find('.btn').addClass('disabled').attr('disabled', 'disabled');
        $(this).parents('form').submit();
    });
});