import json
import vk_api

API_VERSION = '5.103'
APP_ID = 7238281

# 2FA обработчик
def auth_handler():
    key = input('Введите код двухфакторной аутентификации: ')
    remember_device = True
    return key, remember_device

# captcha обработчик
def captcha_handler(captcha):
    key = input(f"Введите капчу ({captcha.get_url()}): ").strip()
    return captcha.try_again(key)

def get_albums(vk, vk_tools):
    albums = []
    for album in vk.photos.getAlbums(need_system=1)['items']:
        photos = vk_tools.get_all(values={'album_id': album['id'], 'photo_sizes': 1}, method='photos.get',
                                  max_count=1000)
        print(photos)
        albums.append({'name': album['title'], 'photos': [p['sizes'][-1]['url'] for p in photos['items']]})
    return albums

def login_vk(login=None, password=None):
    vk_session = vk_api.VkApi(login, password, captcha_handler=captcha_handler, app_id=APP_ID, 
                              api_version=API_VERSION, auth_handler=auth_handler, scope = 'PHOTOS')
    vk_session.auth(token_only=True, reauth=True)
    return vk_session.get_api()

def main():
    inputLogin = input('Login: ')
    inputPassword = input('Password: ')
    vk = login_vk(login=inputLogin, password=inputPassword)
    vk_tools = vk_api.VkTools(vk)
    albums = get_albums(vk, vk_tools)
    print(json.dumps(albums))
    

if __name__ == '__main__':
    main()
