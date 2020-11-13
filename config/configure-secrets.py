#!/usr/bin/python

import sys
import yaml
import subprocess

def usage():
  print('''
    NAME
        configure-secrets - Create or delete platform deployment secrets from FILE

    SYNOPSIS
        configure-secrets <ACTION> FILE

    ACTION
        -s, set     Set secrets
        -d, del     Delete secrets

    FILE
        File containing secrets to be provisioned/removed
  ''')

# Parse secrets file
def parse(fname):
  print '\n>>> Parsing secrets file'
  with open(fname, 'r') as stream:
    secrets = {}
    try:
      secrets = yaml.safe_load(stream)
    except yaml.YAMLError as exc:
      print(exc)
      print('ERROR: failed to parse yaml file')
      exit(1)
    return secrets

# Add provided secrets 
def add(secrets):
  print('\n>>> Setting secrets')
  if not bool(secrets):
    print('no secrets to add')
    return

  for secret, fields in secrets.items():
    if not bool(fields):
      print('skipping secret with no fields: ' + secret)
      continue

    entries = ''
    for field, value in fields.items():
      entries += ' --from-literal=' + field + '=' + value

    subprocess.call('kubectl create secret generic ' + secret + entries, shell=True)

# Remove provided secrets
def remove(secrets):
  print('\n>>> Removing secrets')
  if not bool(secrets):
    print('no secrets to remove')
    return

  for secret, fields in secrets.items():
    subprocess.call('kubectl delete secret ' + secret, shell=True)

# Parse arguments
argCount = len(sys.argv)
if argCount != 3:
  print('ERROR: invalid number of args')
  usage()
  sys.exit(1)
action = sys.argv[1]
fname = sys.argv[2]

# Run command
if (action == '-s' or action == 'set'):
  secrets = parse(fname)
  remove(secrets)
  add(secrets)
elif (action == '-d' or action == 'del'):
  secrets = parse(fname)
  remove(secrets)
else:
  print('ERROR: invalid action')
  usage()
  sys.exit(1)

print('')


