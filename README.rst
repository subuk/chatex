
Live anonymous image board with websockets and html5 notifications.

For websockets fallback work properly, you should configure flash socket policy server.

Nginx configuration example:

::

    server {
        listen   843;
        root /usr/share/nginx/html/flash;

        location / {
            rewrite ^(.*)$ /flash-policy.xml;
        }

        location /flash-policy.xml {
            root /usr/share/nginx/html/flash;
        }

        error_page 400 /flash-policy.xml;

    }

flash-policy.xml:

::

    <cross-domain-policy>
        <allow-access-from domain="board.example.org" to-ports="*"/>
    </cross-domain-policy>
