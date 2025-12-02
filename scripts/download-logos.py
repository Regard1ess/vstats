#!/usr/bin/env python3
"""
下载常见厂商和发行版的logo资源
使用simple-icons CDN作为主要来源
"""

import os
import sys
from pathlib import Path
from typing import Dict, List, Tuple
from urllib.request import urlopen, Request
from urllib.error import URLError, HTTPError

# 脚本目录
SCRIPT_DIR = Path(__file__).parent.absolute()
PROJECT_ROOT = SCRIPT_DIR.parent
LOGOS_DIR = PROJECT_ROOT / "web" / "public" / "logos"
PROVIDERS_DIR = LOGOS_DIR / "providers"
DISTRIBUTIONS_DIR = LOGOS_DIR / "distributions"

# 厂商列表 (name, official_logo_url, format)
# 使用官方原始logo URL，优先使用Wikimedia Commons（最可靠）
PROVIDERS: List[Tuple[str, str, str]] = [
    ("AWS", "https://upload.wikimedia.org/wikipedia/commons/9/93/Amazon_Web_Services_Logo.svg", "svg"),
    ("Alibaba Cloud", "https://upload.wikimedia.org/wikipedia/commons/1/17/Alibaba_Cloud_logo.svg", "svg"),
    ("Tencent Cloud", "https://upload.wikimedia.org/wikipedia/commons/9/9e/Tencent_Cloud_logo.svg", "svg"),
    ("Vultr", "https://upload.wikimedia.org/wikipedia/commons/0/0e/Vultr_logo.svg", "svg"),
    ("DigitalOcean", "https://upload.wikimedia.org/wikipedia/commons/f/ff/DigitalOcean_logo.svg", "svg"),
    ("Linode", "https://upload.wikimedia.org/wikipedia/commons/5/5a/Linode_logo.svg", "svg"),
    ("Bandwagon", "https://bandwagonhost.com/images/logo.png", "png"),
    ("Huawei Cloud", "https://upload.wikimedia.org/wikipedia/commons/0/0e/Huawei_Standard_logo.svg", "svg"),
    ("Google Cloud", "https://upload.wikimedia.org/wikipedia/commons/5/5e/Google_Cloud_logo.svg", "svg"),
    ("Azure", "https://upload.wikimedia.org/wikipedia/commons/a/a8/Microsoft_Azure_Logo.svg", "svg"),
    ("Oracle Cloud", "https://upload.wikimedia.org/wikipedia/commons/5/50/Oracle_logo.svg", "svg"),
    ("IBM Cloud", "https://upload.wikimedia.org/wikipedia/commons/5/51/IBM_logo.svg", "svg"),
    ("Cloudflare", "https://upload.wikimedia.org/wikipedia/commons/9/94/Cloudflare_Logo.svg", "svg"),
    ("OVH", "https://upload.wikimedia.org/wikipedia/commons/9/9e/OVH_logo.svg", "svg"),
    ("Hetzner", "https://upload.wikimedia.org/wikipedia/commons/0/0a/Hetzner_Logo.svg", "svg"),
    ("Scaleway", "https://www.scaleway.com/assets/images/logo-scaleway.svg", "svg"),
    ("Contabo", "https://contabo.com/assets/img/logo.svg", "svg"),
    ("Kamatera", "https://www.kamatera.com/wp-content/themes/kamatera/images/logo.svg", "svg"),
    ("Rackspace", "https://upload.wikimedia.org/wikipedia/commons/7/7a/Rackspace_logo.svg", "svg"),
    ("Joyent", "https://upload.wikimedia.org/wikipedia/commons/9/9a/Joyent_logo.svg", "svg"),
]

# 发行版列表 (name, official_logo_url, format)
# 使用官方原始logo URL，优先使用Wikimedia Commons或官方CDN
DISTRIBUTIONS: List[Tuple[str, str, str]] = [
    ("Ubuntu", "https://upload.wikimedia.org/wikipedia/commons/a/ab/Ubuntu_logo.svg", "svg"),
    ("Debian", "https://www.debian.org/logos/openlogo-nd.svg", "svg"),
    ("CentOS", "https://upload.wikimedia.org/wikipedia/commons/b/bf/Centos-logo-light.svg", "svg"),
    ("Rocky Linux", "https://rockylinux.org/images/logo.svg", "svg"),
    ("AlmaLinux", "https://almalinux.org/images/logo.svg", "svg"),
    ("Fedora", "https://upload.wikimedia.org/wikipedia/commons/3/3f/Fedora_logo.svg", "svg"),
    ("Red Hat", "https://upload.wikimedia.org/wikipedia/commons/a/a8/Red_Hat_logo.svg", "svg"),
    ("SUSE", "https://upload.wikimedia.org/wikipedia/commons/b/b5/Suse_Logo.svg", "svg"),
    ("openSUSE", "https://upload.wikimedia.org/wikipedia/commons/7/76/OpenSUSE_Logo.svg", "svg"),
    ("Arch Linux", "https://upload.wikimedia.org/wikipedia/commons/a/a5/Archlinux-icon-crystal-64.svg", "svg"),
    ("Manjaro", "https://upload.wikimedia.org/wikipedia/commons/3/3e/Manjaro-logo.svg", "svg"),
    ("Gentoo", "https://upload.wikimedia.org/wikipedia/commons/4/48/Gentoo_Linux_logo_matte.svg", "svg"),
    ("Slackware", "https://upload.wikimedia.org/wikipedia/commons/2/22/Slackware_logo_from_the_official_Slackware_site.svg", "svg"),
    ("Mint", "https://upload.wikimedia.org/wikipedia/commons/3/3f/Linux_Mint_logo_without_wordmark.svg", "svg"),
    ("Elementary OS", "https://upload.wikimedia.org/wikipedia/commons/8/8a/Elementary_OS_logo.svg", "svg"),
    ("Pop!_OS", "https://pop.system76.com/images/pop-logo.svg", "svg"),
    ("Kali Linux", "https://upload.wikimedia.org/wikipedia/commons/6/6b/Kali_Linux_2.0_wordmark.svg", "svg"),
    ("Parrot OS", "https://upload.wikimedia.org/wikipedia/commons/3/3e/Parrot_Security_Logo.svg", "svg"),
    ("Windows", "https://upload.wikimedia.org/wikipedia/commons/4/44/Microsoft_logo.svg", "svg"),
    ("macOS", "https://upload.wikimedia.org/wikipedia/commons/f/fa/Apple_logo_black.svg", "svg"),
    ("FreeBSD", "https://upload.wikimedia.org/wikipedia/commons/d/dc/Freebsd_logo.svg", "svg"),
    ("OpenBSD", "https://upload.wikimedia.org/wikipedia/commons/8/8b/OpenBSD_Logo_-_Cartoon_Puffy.svg", "svg"),
    ("NetBSD", "https://upload.wikimedia.org/wikipedia/commons/0/0c/NetBSD.svg", "svg"),
    ("DragonFly BSD", "https://upload.wikimedia.org/wikipedia/commons/7/7a/Dragonfly_BSD_logo.svg", "svg"),
]

# 不再使用simple-icons，改用官方logo


def sanitize_filename(name: str) -> str:
    """将名称转换为安全的文件名"""
    # 替换特殊字符
    name = name.replace(" ", "_").replace("!", "").replace("/", "_")
    name = name.replace("Cloud", "Cloud").replace("Linux", "Linux")
    return name.lower()


def download_logo(url: str, filepath: Path, name: str, file_format: str = "svg") -> bool:
    """下载logo文件"""
    try:
        req = Request(url)
        req.add_header('User-Agent', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36')
        req.add_header('Accept', 'image/svg+xml,image/png,image/*,*/*')
        
        with urlopen(req, timeout=15) as response:
            content_type = response.headers.get('Content-Type', '')
            content = response.read()
            
            # 根据文件格式处理
            if file_format == "svg" or url.endswith('.svg') or 'svg' in content_type.lower():
                try:
                    text_content = content.decode('utf-8')
                    # 确保是有效的SVG
                    if not text_content.strip().startswith('<svg') and not text_content.strip().startswith('<?xml'):
                        print(f"  ⚠️  {name}: 下载的内容不是有效的SVG")
                        return False
                    # 保存为文本
                    filepath.parent.mkdir(parents=True, exist_ok=True)
                    filepath.write_text(text_content, encoding='utf-8')
                except UnicodeDecodeError:
                    print(f"  ⚠️  {name}: 无法解码为UTF-8")
                    return False
            else:
                # PNG或其他二进制格式
                filepath.parent.mkdir(parents=True, exist_ok=True)
                filepath.write_bytes(content)
        
        return True
    except HTTPError as e:
        print(f"  ❌ HTTP错误 {e.code}")
        return False
    except URLError as e:
        print(f"  ❌ URL错误 - {e.reason}")
        return False
    except Exception as e:
        print(f"  ❌ 下载失败 - {e}")
        return False




def download_provider_logos():
    """下载厂商logo"""
    print("\n📦 下载厂商logo...")
    print("=" * 50)
    
    success_count = 0
    total_count = len(PROVIDERS)
    
    for name, url, file_format in PROVIDERS:
        filename = sanitize_filename(name)
        filepath = PROVIDERS_DIR / f"{filename}.{file_format}"
        
        print(f"📥 {name}...", end=" ")
        
        if download_logo(url, filepath, name, file_format):
            print("✅")
            success_count += 1
        else:
            print("❌")
    
    print(f"\n✅ 完成: {success_count}/{total_count} 个厂商logo下载成功")
    return success_count == total_count


def download_distribution_logos():
    """下载发行版logo"""
    print("\n🐧 下载发行版logo...")
    print("=" * 50)
    
    success_count = 0
    total_count = len(DISTRIBUTIONS)
    
    for name, url, file_format in DISTRIBUTIONS:
        filename = sanitize_filename(name)
        filepath = DISTRIBUTIONS_DIR / f"{filename}.{file_format}"
        
        print(f"📥 {name}...", end=" ")
        
        if download_logo(url, filepath, name, file_format):
            print("✅")
            success_count += 1
        else:
            print("❌")
    
    print(f"\n✅ 完成: {success_count}/{total_count} 个发行版logo下载成功")
    return success_count == total_count


def create_index_files():
    """创建索引文件，列出所有可用的logo"""
    import json
    
    # 创建厂商索引
    providers_index = {
        "providers": [
            {
                "name": name,
                "filename": f"{sanitize_filename(name)}.{file_format}",
                "path": f"logos/providers/{sanitize_filename(name)}.{file_format}",
                "format": file_format
            }
            for name, url, file_format in PROVIDERS
        ]
    }
    
    # 创建发行版索引
    distributions_index = {
        "distributions": [
            {
                "name": name,
                "filename": f"{sanitize_filename(name)}.{file_format}",
                "path": f"logos/distributions/{sanitize_filename(name)}.{file_format}",
                "format": file_format
            }
            for name, url, file_format in DISTRIBUTIONS
        ]
    }
    
    index_file = LOGOS_DIR / "index.json"
    index_data = {**providers_index, **distributions_index}
    index_file.write_text(json.dumps(index_data, indent=2, ensure_ascii=False), encoding='utf-8')
    print(f"\n📄 索引文件已创建: {index_file}")


def main():
    """主函数"""
    print("🚀 vStats Logo下载工具")
    print("=" * 50)
    
    # 检查目录
    if not LOGOS_DIR.exists():
        LOGOS_DIR.mkdir(parents=True, exist_ok=True)
    
    # 下载logo
    providers_ok = download_provider_logos()
    distributions_ok = download_distribution_logos()
    
    # 创建索引文件
    create_index_files()
    
    # 总结
    print("\n" + "=" * 50)
    if providers_ok and distributions_ok:
        print("✅ 所有logo下载完成！")
        return 0
    else:
        print("⚠️  部分logo下载失败，请检查网络连接或logo源")
        return 1


if __name__ == "__main__":
    sys.exit(main())

