# Carbon Dioxide Emission Estimation

There are three primary CO2 emission calculation methods.

"Out-of-the-box" SusQL reports an estimated CO2 emission value for all measured workloads using the `static` method:

## `static` Method
- This `static` method uses a static "carbon intensity value" as a coefficient to calculate grams of CO2 emitted.
  This calculation method is used when the `CARBON-METHOD` `ConfigMap` value is set to `static`.

### `static` Method `ConfigMap` Configurable items
  - `CARBON-METHOD` - The `static` method is enabled when this is set to `static`.
  - `CARBON-INTENSITY` - Carbon intensity value. A coefficient used to convert Joules to grams of CO2 per Joule. The unit definition is grams of CO2 per Joule.
    The default carbon intensity value is based on [US EPA](https://www.epa.gov/energy/greenhouse-gases-equivalencies-calculator-calculations-and-references) data.

## `simpledynamic` Method
- The `simpledynamic` method periodically queries the carbon intensity value for a user specified location to provide a more accurate estimation of CO2 emission.
  This calculation method is used when the `CARBON-METHOD` `ConfigMap` value is set to `simpledynamic`.

### `simpledynamic` Method `ConfigMap` Configurable Items
  - `CARBON-METHOD` - The `simpledynamic` method is enabled when this is set to `simpledynamic`.
  - `CARBON-INTENSITY` - This value is set automatically with the `simpledynamic` method. User specified values will be overwritten. The unit definition is grams of CO2 per Joule.
  - `CARBON-INTENSITY-URL` - Specifies a web API that returns carbon intensity. (The default value works as of the date of this writing, but may need to be modified in the future.)
  - `CARBON-LOCATION` - The location identifiers are defined by the API provider, but can be selected by the user. (The default value works as of the date of this writing.)
  - `CARBON-QUERY-RATE` - Interval in seconds at which the carbon intensity is queried. The data available from the source is updated less than hourly, so an interval of greater than 3600 seconds is recommended.
  - `CARBON-QUERY-FILTER` - When the return value is embedded in a JSON object, this specification enables the extraction of the data. The default value matches the default provider.
  - `CARBON-QUERY-CONV-2J` - If the carbon data provider does not provide data in the standard "grams of CO2 per Joule" then this factor can be specified to normalize the units displayed. The default value ensures that the default provider data is in the correct unit.

## `casdk` Method
- The `casdk` offers integration with the Green Software Foundation's [Carbon Aware SDK](https://github.com/Green-Software-Foundation/carbon-aware-sdk).
  This calculation method is used when the `CARBON-METHOD` `ConfigMap` value is set to `casdk`.
  The user is required to first prepare a local instance of the Carbon Aware SDK that is configured to support carbon intensity queries.

### Configuring and installing Carbon Aware SDK
- Following guidance in https://github.com/Green-Software-Foundation/carbon-aware-sdk/blob/dev/casdk-docs/docs/overview/enablement.md,
the Carbon Aware SDK can be easily installed on a Kubernetes cluster such as OpenShift:
- Preparation: clone the repository and edit `helm-chart/values.yaml` as needed to reflect private password, configuration, etc.
(Useful configuration tips available at https://github.com/Green-Software-Foundation/carbon-aware-sdk/blob/dev/casdk-docs/docs/tutorial-extras/configuration.md )

```
git clone git@github.com:Green-Software-Foundation/carbon-aware-sdk.git
vi helm-chart/values.yaml
```
- Preparation: required software and permission
  - Ensure that `helm`, and `kubectl` (or `oc`) are installed
  - Ensure that CLI user is logged in to cluster with sufficient permissions
- Perform installation
```
cd carbon-aware-sdk
helm upgrade --install --wait carbon-aware-sdk helm-chart --create-namespace gsf
oc expose svc/webapi -n gsf
oc get routeso
```
Note the value reported for "HOST/PORT". This will be used in the next configuration step.
- Update `susql-config.yaml`
Update the following items in the `susql-config.yaml` file:
```
  CARBON-METHOD: "casdk"
  CARBON-INTENSITY-URL: "http://<HOST/PORT-VALUE>/carbonn-intensity/latest?zone=%s"
  CARBON-LOCATION: "<YOUR-LOCATION>"
```
Apply the updated `susql-config.yaml` file:
```
oc apply -f susql-config.yaml -n <SUSQL-OPERATOR-NAMESPACE>
```
You are now ready to install and use the SusQL operator.
If the SusQL Operator is alreay installed, then restart the control pod.


### `casdk` Method `ConfigMap` Configurable Items
  - `CARBON-METHOD` - The `casdk` method is enabled when this is set to `casdk`.
  - `CARBON-INTENSITY` - This value is set automatically with the `casdk` method. User specified values will be overwritten. The unit definition is grams of CO2 per Joule.
  - `CARBON-INTENSITY-URL` - Specifies the web API that returns carbon intensity. This is assumed to be a GSF Carbon Aware SDK prepared and configured by the user.
  - `CARBON-LOCATION` - The location identifiers are defined by the API provider, and should be specified by the user.
  - `CARBON-QUERY-RATE` - Interval in seconds at which the carbon intensity is queried.
  - `CARBON-QUERY-CONV-2J` - The default values converts "grams of CO2 per KWH" (Carbon Aware SDK standard) to "grams of CO2 per Joule".
  - `moremoremore` - more more more more.

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

