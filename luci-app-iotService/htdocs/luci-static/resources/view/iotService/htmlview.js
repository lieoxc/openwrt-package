'use strict';
'require uci';
'require view';
'require ui';
return view.extend({
	load: function () {
		
	},
	render: function (data) {
		var body = E([
		]);
		var listContainer = E('div');
		var list = E('ul');
		list.appendChild(E('button', {
			'class': 'btn cbi-button-action',
			'click': ui.createHandlerFn(this, function(ev) {
				window.open('http://192.168.10.1:9999', '_blank')
			})
		}, ['%:Go to Iot Service%']));
		listContainer.appendChild(list);
		body.appendChild(listContainer);
		return body;
	}
});