$(function(){

	$.fn.showModal = function(){
		$(this).modal('setting', {blurring:true}).modal('show')
	}

	// object gets an object describing all the input values
	// for a form.
	$.fn.object = function(){
		var obj = {}
		$(this).find('input,textarea').each(function(){
			var $this = $(this),
				name = $this.attr('name'),
				value = $this.val()
			if (name && value)
				obj[name] = value
		})
		return obj
	}

	// api is a wrapper for $.ajax that does common things.
	$.api = function(options){
		if (options.form) {
			options.type = options.form.attr('method') || 'get'
			options.url = options.form.attr('action') || ""
			options.data = JSON.stringify(options.form.object())
			options.dataType = 'json'
			options.contentType = 'application/json'
		}
		options._error = options.error
		options.error = function(response){
			console.warn(response)
			var message = response.responseText || response.statusText || "An unknown error occurred"
			if (response.responseJSON && response.responseJSON.error) {
				message = response.responseJSON.error
			}
			options._error = options._error || function(message){
				// global error handler
				var errEl = $('.ui.global.error.message')
				if (errEl.length == 0) {
					errEl = $("<div>", {class:"ui global error message container"}).insertAfter($(".topnav"))
				}
				errEl.text(message)
			}
			options._error(message, response)
		}
		$.ajax(options)
	}

	// inject sets the text or value of page elements to the
	// data represented in the argument.
	$.fn.inject = function(data){
		var $this = $(this)
		for (var k in data) {
			if (!data.hasOwnProperty(k)) continue
			$this.find('[data-field="'+k+'"]').each(function(){
				var $that = $(this)
				switch ($that[0].tagName) {
					case "INPUT":
						$that.val(data[k])
						break
					default:
						$that.text(data[k])
				}
			})
		}
	}

	// data-trigger-modal events
	$('[data-trigger-modal]').on('click', function(e){
		e.preventDefault()
		var modal = $(this).attr('data-trigger-modal')
		$('[data-modal="'+modal+'"]').showModal()
	})

	// re-insert all deferred scripts that may already
	// have been put into the page
	$('script[type="deferred"]')
		.remove()
		.attr('type', 'text/javascript')
		.appendTo($('body'))

})