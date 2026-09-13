#!/usr/bin/env python3
"""Local QEMU + Mosquitto integration. Firmware must use samples/network.conf.

Requires qemu-system-riscv32, mosquitto, mosquitto_pub, mosquitto_sub on PATH.
Uses only an isolated local broker; does not flash hardware or contact a cloud.
"""
import argparse
import json
from pathlib import Path
import shutil
import socket
import subprocess
import tempfile
import time


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('firmware', type=Path)
    parser.add_argument('--port', type=int, default=1883)
    args = parser.parse_args()
    for name in ('qemu-system-riscv32', 'mosquitto', 'mosquitto_pub', 'mosquitto_sub'):
        if not shutil.which(name):
            raise SystemExit(f'Missing tool: {name}')
    with socket.socket() as check:
        check.bind(('127.0.0.1', args.port))
    root = Path(tempfile.mkdtemp(prefix='zephyr-mqtt-'))
    print(f'Logs: {root}', flush=True)
    config = root / 'mosquitto.conf'
    config.write_text(f'listener {args.port} 127.0.0.1\nallow_anonymous true\npersistence false\n')
    processes = []
    streams = []

    def start(command, log):
        stream = (root / log).open('a')
        streams.append(stream)
        p = subprocess.Popen(command, stdin=subprocess.DEVNULL, stdout=stream, stderr=stream)
        processes.append(p)
        return p

    def broker():
        return start(['mosquitto', '-c', str(config), '-v'], 'broker.log')

    def messages():
        path = root / 'messages.log'
        found = []
        if path.exists():
            for line in path.read_text().splitlines():
                try:
                    topic, payload = line.split(' ', 1)
                    found.append((topic, json.loads(payload)))
                except (ValueError, json.JSONDecodeError):
                    pass
        return found

    def wait_for(predicate, label, timeout=45):
        deadline = time.monotonic() + timeout
        while time.monotonic() < deadline:
            if predicate():
                print('PASS:', label, flush=True)
                return
            time.sleep(.2)
        raise AssertionError(f'Timeout: {label}; inspect {root}')

    def send(suffix, payload):
        subprocess.run(['mosquitto_pub', '-h', '127.0.0.1', '-p', str(args.port),
                        '-q', '1', '-t', f'devices/{device_id}/{suffix}',
                        '-m', json.dumps(payload)], check=True, timeout=10)

    try:
        server = broker()
        time.sleep(.5)
        if server.poll() is not None:
            raise AssertionError('Broker did not start')
        start(['mosquitto_sub', '-h', '127.0.0.1', '-p', str(args.port),
               '-q', '1', '-v', '-t', 'devices/#'], 'messages.log')
        start(['qemu-system-riscv32', '-machine', 'virt', '-m', '256M',
               '-nographic', '-bios', 'none', '-global', 'virtio-mmio.force-legacy=false', '-netdev', 'user,id=n1',
               '-device', 'virtio-net-device,bus=virtio-mmio-bus.0,netdev=n1,mac=52:54:00:12:34:56',
               '-kernel', str(args.firmware.resolve())], 'device.log')
        wait_for(lambda: any(t.endswith('/presence') and p.get('online') for t, p in messages()),
                 'MQTT connected and retained presence')
        device_id = next(t.split('/')[1] for t, p in messages() if t.endswith('/presence') and p.get('online'))
        send('config', {'version': 1, 'collectors': {
            'temperature': {'collection_interval_sec': 1, 'upload_interval_sec': 3},
            'imu': {'collection_interval_sec': 2, 'upload_interval_sec': 5},
            'battery': {'enabled': False}}})
        wait_for(lambda: any(p.get('type') == 'config_applied' for _, p in messages()),
                 'remote config applies without reboot')
        wait_for(lambda: any(t.endswith('/telemetry') and p.get('collector') == 'temperature'
                             and len(p.get('records', [])) >= 3 for t, p in messages()),
                 'collection and upload intervals differ')
        device_id = next(t.split('/')[1] for t, p in messages() if t.endswith('/presence') and p.get('online'))
        send('config', {'version': 1, 'collectors': {'temperature': {'enabled': False}}})
        wait_for(lambda: any(p.get('type') == 'config_rejected' for _, p in messages()),
                 'stale retained version rejected')
        before = sum(t.endswith('/telemetry') for t, _ in messages())
        server.terminate()
        server.wait(timeout=5)
        time.sleep(8)
        server = broker()
        wait_for(lambda: sum(t.endswith('/telemetry') for t, _ in messages()) > before + 1,
                 'broker recovery automatically flushes offline cache', timeout=75)
        batches = [p for t, p in messages() if t.endswith('/telemetry')]
        assert any(p.get('collector') == 'imu' for p in batches), 'independent IMU missing'
        assert all('sequence' in record and 'time_synced' in record
                   for batch in batches for record in batch['records'])
        print(f'PASS: {len(batches)} telemetry batches; two independent collectors', flush=True)
    finally:
        for p in reversed(processes):
            if p.poll() is None:
                p.terminate()
                try:
                    p.wait(timeout=5)
                except subprocess.TimeoutExpired:
                    p.kill()
                    p.wait()
        for stream in streams:
            stream.close()


if __name__ == '__main__':
    main()
