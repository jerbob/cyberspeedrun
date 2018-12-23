from json import load
from time import time
from threading import Thread

from requests_html import HTMLSession


go_endpoint = 'https://go.joincyberdiscovery.com'

with open('solutions.json') as file:
    solutions = load(file)

session = HTMLSession()


def find_csrf(response):
    """Locate the CSRF token, given a response object."""
    return response.html.xpath(
        '//*[@type="hidden"]/@value'
    )[0]


def post_flag(flag, challenge, csrf):
    """Post the flag for a challenge with the given session object."""
    print(f'[+] Flag found: {flag}')
    return session.post(
        f'{go_endpoint}/api',
        data={
            'action': 'flag_attempt',
            'flag': flag,
            'challenge': challenge,
            'csrf': csrf
        }
    )


def solve_challenge(challenge, solution):
    """Solve a challenge, given the name and solution."""
    csrf = find_csrf(
        session.get(f'{go_endpoint}/challenges/{challenge}')
    )
    # Static Flag
    if type(solution) is not dict:
        post_flag(solution, challenge, csrf)
    # Requires POST to XHR route
    else:
        flag = session.post(
            f'{go_endpoint}/challenges/xhr/{challenge}',
            data={**solution, 'csrf': csrf}
        ).json()['attempt']['code']
        post_flag(flag, challenge, csrf)


def main(username):
    """Speedrun Go with the specified username."""
    threads = []
    csrf = find_csrf(session.get(f'{go_endpoint}/register'))
    print('\n[?] Registering...')

    session.post(
        f'{go_endpoint}/register',
        data={'csrf': csrf, 'name': username}
    )
    print(f'[!] Registered as {username}.\n')
    start = time()
    print(f'[!] Started timer.\n')

    for challenge, solution in solutions.items():
        print(f'[?] Attempting challenge {challenge}...')
        thread = Thread(target=solve_challenge, args=(challenge, solution))
        threads.append(thread)
        thread.start()

    print()
    for thread in threads:
        thread.join()

    print(f'\n[!] Completed speedrun in {round(time() - start, 2)}s.')


main('Aperture')
