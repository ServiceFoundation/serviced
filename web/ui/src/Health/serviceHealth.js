/* jshint multistr: true */
(function() {
    'use strict';

    // OK means health check is passing
    const OK = "passed";
    // Failed means health check is responsive, but failing
    const FAILED = "failed";
    // Timeout means health check is non-responsive in the given time
    const TIMEOUT = "timeout";
    // NotRunning means the instance is not running
    const NOT_RUNNING = "not_running";
    // Unknown means the instance hasn't checked in within the provided time
    // limit.
    const UNKNOWN = "unknown";
    // EMERGENCY_SHUTDOWN means instance has been emergency shutdown
    const EMERGENCY_SHUTDOWN = "emergency_shutdown";

    let serviceHealthModule = angular.module('serviceHealth', []);

    // share constants for other packages to use
    serviceHealthModule.value("hcStatus", {
        OK: OK,
        FAILED: FAILED,
        TIMEOUT: TIMEOUT,
        NOT_RUNNING: NOT_RUNNING,
        UNKNOWN: UNKNOWN,
        EMERGENCY_SHUTDOWN: EMERGENCY_SHUTDOWN
    });

    serviceHealthModule.factory("$serviceHealth", ["$translate",
    function($translate){

        var statuses = {};

        // updates health check data for all services
        function update(serviceList) {

            var serviceStatus, instanceStatus, instanceUniqueId, service;

            statuses = {};

            // iterate services healthchecks
            for(var serviceId in serviceList){
                service = serviceList[serviceId];
                serviceStatus = new Status(
                    serviceId,
                    service.name,
                    service.desiredState,
                    service.emergencyShutdown);

                // refresh list of instances
                // TODO - this "if" is a workaround for old servicesFactory
                // services and should be removed along with servicesFactory
                if(service.fetchInstances){
                    service.fetchInstances();
                }

                // if this service has instances, evaluate their health
                service.instances.forEach(instance => {

                    // create a new status rollup for this instance
                    instanceUniqueId = serviceId +"."+ instance.model.InstanceID;
                    instanceStatus = new Status(
                        instanceUniqueId,
                        service.name +" "+ instance.model.InstanceID,
                        service.desiredState,
                        service.emergencyShutdown);

                    // evalute instance healthchecks and roll em up
                    instanceStatus.evaluateHealthChecks(instance.healthChecks);
                    // store resulting status on instance
                    instance.status = instanceStatus;

                    // add this guy's statuses to hash map for easy lookup
                    statuses[instanceUniqueId] = instanceStatus;
                    // add this guy's status to his parent
                    serviceStatus.children.push(instanceStatus);
                });

                // now that this services instances have been evaluated,
                // evaluate the status of this service
                serviceStatus.evaluateChildren();

                statuses[serviceId] = serviceStatus;
            }

            return statuses;
        }

        function evaluate(service, instances){

            instances = instances || [];

            let status;

            status = new Status(
                service.id,
                service.name,
                service.desiredState,
                service.emergencyShutdown);

            // if instances were provided, evaluate their health
            instances.forEach(instance => {
                let instanceUniqueId = service.id +"."+ instance.model.InstanceID;
                let instanceStatus = new Status(
                    instanceUniqueId,
                    service.name +" "+ instance.model.InstanceID,
                    service.desiredState,
                    service.emergencyShutdown);

                // evalute instance healthchecks and roll em up
                instanceStatus.evaluateHealthChecks(instance.healthChecks);
                // store resulting status on instance
                instance.status = instanceStatus;

                // add this guy's status to his parent
                status.children.push(instanceStatus);
            });

            // now that this services instances have been evaluated,
            // evaluate the status of this service
            status.evaluateChildren();

            return status;
        }

        // used by Status to examine children and figure
        // out what the parent's status is
        function StatusRollup(){
            this[OK] = 0;
            this[FAILED] = 0;
            this[NOT_RUNNING] = 0;
            this[UNKNOWN] = 0;
            this.total = 0;
        }
        StatusRollup.prototype = {
            constructor: StatusRollup,

            incOK: function(){
                this.incStatus(OK);
            },
            incFailed: function(){
                this.incStatus(FAILED);
            },
            incNotRunning: function(){
                this.incStatus(NOT_RUNNING);
            },
            incUnknown: function(){
                this.incStatus(UNKNOWN);
            },
            incStatus: function(status){
                if(this[status] !== undefined){
                    this[status]++;
                    this.total++;
                }
            },

            // TODO - use assertion style ie: status.is.ok() or status.any.ok()
            anyFailed: function(){
                return !!this[FAILED];
            },
            allFailed: function(){
                return this.total && this[FAILED] === this.total;
            },
            anyOK: function(){
                return !!this[OK];
            },
            allOK: function(){
                return this.total && this[OK] === this.total;
            },
            anyNotRunning: function(){
                return !!this[NOT_RUNNING];
            },
            allNotRunning: function(){
                return this[NOT_RUNNING] === this.total;
            },
            anyUnknown: function(){
                return !!this[UNKNOWN];
            },
            allUnknown: function(){
                return this.total && this[UNKNOWN] === this.total;
            }
        };

        function Status(id, name, desiredState, emergencyShutdown){
            this.id = id;
            this.name = name;
            this.desiredState = desiredState;
            this.emergencyShutdown = emergencyShutdown;

            this.statusRollup = new StatusRollup();
            this.children = [];

            this.status = null;
            this.description = null;
        }

        Status.prototype = {
            constructor: Status,

            // distill this service's statusRollup into a single value
            evaluateStatus: function(){
                    // stuff is shutdown and emergency shutdown is flagged
                    if(this.emergencyShutdown){
                        this.status = EMERGENCY_SHUTDOWN;
                        this.description = $translate.instant("emergency_shutdown");
                    // if all are ok, up
                    } else if(this.statusRollup.allOK()){
                        this.status = OK;
                        this.description = $translate.instant("passing_health_checks");
                    // if all are failed, fail
                    } else if(this.statusRollup.allFailed()){
                        this.status = FAILED;
                        this.description = $translate.instant("failed");
                    // if all are stopped, stopped
                    } else if(this.statusRollup.allNotRunning()){
                        this.status = NOT_RUNNING;
                        this.description = $translate.instant("container_down");
                    // else some running up some stopped
                    } else {
                    this.status = UNKNOWN;
                    this.description = $translate.instant("missing_health_checks");
                    }
            },

            // roll up child status into this status
            evaluateChildren: function(){
                this.statusRollup = this.children.reduce(function(acc, childStatus){
                    acc.incStatus(childStatus.status);
                    return acc;
                }.bind(this), new StatusRollup());
                this.evaluateStatus();
            },

            // set this status's statusRollup based on healthchecks
            // NOTE - subtly different than evaluateChildren
            evaluateHealthChecks: function(healthChecks){
                for(var name in healthChecks){
                    this.statusRollup.incStatus(healthChecks[name]);
                    this.children.push({
                        name: name,
                        status: healthChecks[name]
                    });
                }
                this.evaluateStatus();
            },

        };

        return {
            update,
            evaluate,
            get: function(id){
                var status = statuses[id];

                // if no status found, return unknown
                if(!status){
                    status = new Status(id, UNKNOWN, 0, false);
                    status.evaluateStatus();
                }

                return status;
            }
        };
    }]);

})();
