#!/bin/bash

if [ -z `getent group qm-hermes` ]; then
	groupadd qm-hermes
fi

if [ -z `getent passwd qm-hermes` ]; then
	useradd qm-hermes -g qm-hermes -s /bin/sh
fi

install --mode=755 --owner=qm-hermes --group=qm-hermes --directory /var/log/qm-hermes

systemctl daemon-reload

#END
