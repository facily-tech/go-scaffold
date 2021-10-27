import http from 'k6/http';
import {group, check} from 'k6';

export let options = {
    vus: 1000,
    duration: '10s',
    summaryTimeUnit: 'ms',
    thresholds: {
        http_req_duration: ['p(95)<500'],
    },
};

const URL = 'http://localhost:8080/v1/quote';
export default function() {
    let body = {
        content: 'quote' + __ITER,
    };
    let params = {
        headers: {
            'Content-Type': 'application/json',
        },
        tags: {
            name: 'create',
        },
    };

    group('api journey', (_) => {
        let create_response = http.post(
            URL,
            JSON.stringify(body),
            params,
        );
        check(create_response, {
            'create status 201': (r) => r.status == 201,
            'create has id?': (r) => r.json().hasOwnProperty('id'),
        });

        let id = create_response.json()['id'];

        params.tags.name = 'find';
        let find_response = http.get(
            URL+ "/" + id,
            null,
            params,
        );
        check(find_response, {
            'find status 200': (r) => r.status == 200,
        });

        params.tags.name = 'update';
        let update_response = http.post(
            URL,
            JSON.stringify({id: id, content: 'updated ' + __ITER}),
            params,
        );
        check(update_response, {
            'update status 201': (r) => r.status == 201,
        });

        params.tags.name = 'delete';
        let delete_response = http.del(
            URL + '/' + id,
            null,
            params,
        );
    });
}