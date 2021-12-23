# create a single user resource
helm install ${usertag} ./helmChartFiles/allInOne --set jupyter.user=${usertag} --set grafana.podLabels.user=${usertag} --set global.user=${usertag} -n user-resource

# disable grafana
helm install ${usertag} -f helmChartFiles/singleUserHelm/values.yaml ./helmChartFiles/singleUserHelm --set jupyter.user=${usertag} --set grafana.enabled=false --set global.user=${usertag} -n user-resource

# disable jupyter
helm install ${usertag} -f helmChartFiles/singleUserHelm/values.yaml ./helmChartFiles/singleUserHelm --set jupyter.enabled=false --set grafana.podLabels.user=${usertag} --set global.user=${usertag} -n user-resource


# install jupyter
helm install ${usertag}-jupyter --set user=$usertag --set jupyter.user=${usertag}-jupyter ./helmChartFiles/jupyterChart -n user-resource

# install grafana
helm install ${usertag}-grafana --set user=$usertag --set grafana.podLabels.user=${usertag}-grafana ./helmChartFiles/grafanaChart -n user-resource