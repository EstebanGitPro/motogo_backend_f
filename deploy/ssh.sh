#!/bin/bash
# ============================================
# SSH Connection - MotoGo VPS (Contabo)
# ============================================
# IP:     161.97.142.2
# OS:     Ubuntu 24.04
# Plan:   Contabo Cloud VPS 20
# ============================================
#
# Primera vez → conectar como root:
#   ./ssh.sh
#
# Después de crear usuario motogo:
#   ./ssh.sh motogo
#
# ============================================

VPS_IP="161.97.142.2"
USER="${1:-root}"

echo "🔌 Conectando a MotoGo VPS..."
echo "   Usuario: $USER"
echo "   IP:      $VPS_IP"
echo ""

ssh -o ConnectTimeout=10 \
    -o ServerAliveInterval=60 \
    -o ServerAliveCountMax=3 \
    "$USER@$VPS_IP"
