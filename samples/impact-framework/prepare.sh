apt update -y
apt upgrade -y
apt upgrade npm vim lsb-release -y
cd /root

cat <<EOF > demo.yaml
name: basic-demo
description:
tags:
initialize:
  plugins:
    double-a-value:
      path: 'builtin'
      method: Coefficient
      config:
        input-parameter: "cpu-utilization"
        coefficient: 2
        output-parameter: "cpu-utilization-doubled"

tree:
  children:
    child-0:
      defaults:
        cpu/thermal-design-power: 100
      pipeline:
        observe:
        regroup:
        compute:
          - double-a-value
      inputs:
        - timestamp: 2023-07-06T00:00
          duration: 1
          cpu/utilization: 20
        - timestamp: 2023-07-06T00:01
          duration: 1
          cpu/utilization: 80
EOF
git clone https://github.com/Green-Software-Foundation/if
cd if
npm install
npm run if-run -- --manifest  manifests/examples/builtins/sum/success.yml
