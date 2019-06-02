#############################
##
## Condor job specification
##
##############################

# This defines what job universe we want the job to run in.
# 'vanilla' is the simplest option for basic command execution.
# Other universes that exist include 'standard', 'grid', 'java' and 'mpi'.
universe        = vanilla

# This defines the path of the executable we want to run.
executable      = doc_train.sh

# This specifies where data sent to STDOUT by the executable should be
# directed to.
#
# The Condor system can perform variable substitution in job specifications;
# the $(Process) string below will be replaced with the job's Process value.
# If we submit multiple jobs from this single specification (we do, as you
# will see later) then the Process value will be incremented for each job.
# For example, if we submit 100 jobs, then each job will have a different
# Process value, numbered in sequence from 0 to 99.
#
# If we were to instruct every job to redirect STDOUT to the same file, then
# data would be lost as each job would overwrite the same file in an
# uncontrolled manner.  Thus, we direct STDOUT for each job to a uniquely
# named file.
output          = out/logs/labeller.$(Process).out

# As above, but for STDERR.
error           = out/err/labeller.$(Process).out

# Condor can write a time-ordered log of events related to this job-set
# to a file we specify.  This specifies where that file should be written.
log             = out/events.log

# Set environment variables
environment = "TRANSPORT_DB_USERNAME={}"

# This specifies what commandline arguments should be passed to the executable.
arguments       = -p random_forest

# This specifies that the specification, as parsed up to this point, should be
# submitted 5 times.  (If the number is omitted, the number '1' is assumed.)
queue 150